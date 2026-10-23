// Package gallery owns all communication with the Divoom gallery servers.
package gallery

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"divoom-gallery-client/internal/animation"
	"divoom-gallery-client/internal/httpx"
	"divoom-gallery-client/internal/media"
)

const galleryAgent = "Aurabox/3.1.10 (iPad; iOS 14.8; Scale/2.00)"
const maxResponseSize = 2 * 1024 * 1024

// Client fetches gallery data. It owns authentication and session renewal,
// without knowing about local files, playback history, rotation, or web pages.
type Client struct {
	mu              sync.Mutex
	api, files      *httpx.Client
	email, password string
	token, userID   json.RawMessage
}

func New(options Options) *Client {
	apiURL, fileURL := options.APIURL, options.FileURL
	if apiURL == "" {
		apiURL = "https://appin.divoom-gz.com"
	}
	if fileURL == "" {
		fileURL = "https://fin.divoom-gz.com"
	}
	newHTTPClient := func(baseURL string) *httpx.Client {
		return httpx.NewClient(httpx.Options{
			BaseURL: baseURL, Timeout: 30 * time.Second,
			Headers: map[string]string{"User-Agent": galleryAgent},
			Service: "gallery", Metrics: options.Metrics,
		})
	}
	return &Client{api: newHTTPClient(apiURL), files: newHTTPClient(fileURL), email: options.Email, password: options.Password}
}

func (c *Client) FetchCategoryPage(ctx context.Context, category, page int) (Page, error) {
	if category < 0 || category > 255 || page < 1 {
		return Page{}, errors.New("invalid category or page")
	}
	return c.fetchPage(ctx, "/GetCategoryFileListV2", map[string]any{
		"StartNum": (page-1)*PageSize + 1, "EndNum": page * PageSize,
		"Classify": category, "FileSize": 1, "FileType": 1,
		"FileSort": 0, "Version": 19, "RefreshIndex": 0,
	})
}

func (c *Client) FetchAlbumPage(ctx context.Context, albumID int64, page int) (Page, error) {
	if albumID <= 0 || page < 1 {
		return Page{}, errors.New("invalid album or page")
	}
	return c.fetchPage(ctx, "/Discover/GetAlbumImageListV3", map[string]any{
		"AlbumId": albumID, "StartNum": (page-1)*PageSize + 1, "EndNum": page * PageSize,
		"FileSize": 1, "FileType": 1, "FileSort": 0,
		"Version": 19, "Classify": 0, "RefreshIndex": 0,
	})
}

func (c *Client) fetchPage(ctx context.Context, path string, payload map[string]any) (Page, error) {
	var result pageResponse
	if err := c.request(ctx, path, payload, &result); err != nil {
		return Page{}, err
	}
	if result.Items == nil {
		return Page{}, errors.New("gallery response omitted FileList")
	}
	items := make([]Item, 0, len(result.Items))
	for _, item := range result.Items {
		data, err := c.FetchItem(ctx, item)
		if err != nil {
			return Page{}, err
		}
		// FileType describes the gallery category, not the binary encoding.
		// Only downloaded artwork exposes the actual format byte.
		if len(data) > 0 && animation.SupportsFormat(data[0]) {
			items = append(items, item)
		}
	}
	return Page{Items: items, more: len(result.Items) > 0}, nil
}

// FetchItem returns the original artwork bytes for an item returned by a page.
func (c *Client) FetchItem(ctx context.Context, item Item) ([]byte, error) {
	p, err := url.Parse(item.FileID)
	if err != nil || item.FileID == "" || p.IsAbs() || p.Host != "" || p.RawQuery != "" || p.Fragment != "" || strings.Contains(p.Path, "..") || strings.Contains(p.Path, "\\") {
		return nil, errors.New("invalid gallery file path")
	}
	return c.files.Fetch(ctx, httpx.Request{
		Method: "GET", Path: "/" + strings.TrimLeft(p.EscapedPath(), "/"),
		MetricPath: "/artwork", Limit: media.MaxArtworkSize,
	})
}

func (c *Client) ListAlbums(ctx context.Context) ([]Album, error) {
	albums := []Album{}
	seen := make(map[int64]bool)
	// Collect every page for the existing album selectors. A short page alone
	// does not establish the end of a Divoom source; stop at an empty page.
	for page := 1; page <= 10000; page++ {
		var result albumResponse
		if err := c.request(ctx, "/Discover/GetAlbumListV3", map[string]any{
			"StartNum": (page-1)*PageSize + 1, "EndNum": page * PageSize,
			"FileSize": 1, "FileSort": 0,
		}, &result); err != nil {
			return nil, err
		}
		if result.Albums == nil {
			return nil, errors.New("album response omitted AlbumList")
		}
		if len(result.Albums) == 0 {
			return albums, nil
		}
		previousCount := len(albums)
		for _, a := range result.Albums {
			id, err := strconv.ParseInt(strings.Trim(string(a.ID), `"`), 10, 64)
			if err == nil && id > 0 && !seen[id] {
				seen[id] = true
				albums = append(albums, Album{ID: id, Name: a.Name})
			}
		}
		if len(albums) == previousCount {
			return nil, errors.New("album list pagination made no progress")
		}
	}
	return nil, errors.New("album list exceeds page limit")
}

func (c *Client) request(ctx context.Context, path string, payload map[string]any, result any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	authenticated := c.email != "" || c.password != ""
	if authenticated && len(c.token) == 0 {
		if err := c.login(ctx); err != nil {
			return err
		}
	}
	for attempt := 0; attempt < 2; attempt++ {
		if authenticated {
			payload["Token"], payload["UserId"] = c.token, c.userID
		}
		err := c.requestOnce(ctx, path, payload, result)
		var failure *APIError
		if !errors.As(err, &failure) || failure.Code != 11 || !authenticated {
			return err
		}
		c.token, c.userID = nil, nil
		if attempt == 1 {
			return err
		}
		if err := c.login(ctx); err != nil {
			return err
		}
	}
	return errors.New("gallery session renewal failed")
}

func (c *Client) login(ctx context.Context) error {
	hash := md5.Sum([]byte(c.password))
	var result loginResponse
	if err := c.requestOnce(ctx, "/UserLogin", map[string]any{"Email": c.email, "Password": hex.EncodeToString(hash[:])}, &result); err != nil {
		return err
	}
	if len(result.Token) == 0 || string(result.Token) == "null" || string(result.Token) == `""` || len(result.UserID) == 0 || string(result.UserID) == "null" {
		return errors.New("gallery login did not return a session")
	}
	c.token, c.userID = result.Token, result.UserID
	return nil
}

func (c *Client) requestOnce(ctx context.Context, path string, payload map[string]any, result any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return errors.New("invalid gallery request")
	}
	data, err := c.api.Fetch(ctx, httpx.Request{
		Method: "POST", Path: path, Body: bytes.NewReader(body),
		Headers: map[string]string{"Content-Type": "application/json"}, Limit: maxResponseSize,
	})
	if err != nil {
		return err
	}
	var status apiResponse
	if json.Unmarshal(data, &status) != nil {
		return errors.New("invalid gallery response JSON")
	}
	if status.Code != nil && *status.Code != 0 {
		return &APIError{Code: *status.Code}
	}
	if json.Unmarshal(data, result) != nil {
		return errors.New("invalid gallery response fields")
	}
	return nil
}
