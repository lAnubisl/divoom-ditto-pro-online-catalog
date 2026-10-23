package gallery

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListsNeverDownloadArtwork(t *testing.T) {
	files := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/GetCategoryFileListV2", "/Discover/GetAlbumImageListV3":
			fmt.Fprint(w, `{"ReturnCode":0,"FileList":[{"GalleryId":1,"FileId":"artwork","FileType":1,"FileSize":1}]}`)
		default:
			files++
			http.Error(w, "artwork server unavailable", 503)
		}
	}))
	defer server.Close()
	client := New(Options{APIURL: server.URL, FileURL: server.URL})
	for _, fetch := range []func() (Page, error){
		func() (Page, error) { return client.FetchCategoryPage(context.Background(), 0, 1) },
		func() (Page, error) { return client.FetchAlbumPage(context.Background(), 34, 1) },
	} {
		page, err := fetch()
		if err != nil || len(page.Items) != 1 || !page.HasMore() {
			t.Fatalf("metadata listing failed: page=%+v err=%v", page, err)
		}
	}
	if files != 0 {
		t.Fatalf("listing downloaded artwork %d times", files)
	}
}
