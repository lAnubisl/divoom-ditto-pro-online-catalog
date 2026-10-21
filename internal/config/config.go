// Package config is the sole reader of process environment variables.
// It validates startup config and supplies values to application components.
package config

import (
	"errors"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func LoadGalleryAccount() (GalleryAccount, error) {
	a := GalleryAccount{Email: strings.TrimSpace(os.Getenv("DIVOOM_GALLERY_EMAIL")), Password: os.Getenv("DIVOOM_GALLERY_PASSWORD")}
	if a.Email == "" || a.Password == "" {
		return a, errors.New("set DIVOOM_GALLERY_EMAIL and DIVOOM_GALLERY_PASSWORD")
	}
	return a, nil
}

// LoadDeviceAPI reads device API settings from environment variables.
func LoadDeviceAPI() (DeviceAPI, error) {
	d := DeviceAPI{BaseURL: os.Getenv("DIVOOM_API_BASE_URL"), APIKey: os.Getenv("DIVOOM_API_KEY")}
	u, err := url.Parse(strings.TrimSpace(d.BaseURL))
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return d, errors.New("set DIVOOM_API_BASE_URL to an HTTP(S) server URL without credentials, query, or fragment")
	}
	d.BaseURL = strings.TrimRight(u.String(), "/")
	if d.APIKey == "" || strings.ContainsAny(d.APIKey, "\r\n") {
		return d, errors.New("set DIVOOM_API_KEY to a nonempty key without line breaks")
	}
	return d, nil
}

func loadWebPort() (string, error) {
	port := os.Getenv("DIVOOM_WEB_PORT")
	if port == "" {
		return "8090", nil
	}
	number, err := strconv.Atoi(port)
	if err != nil || number < 1 || number > 65535 {
		return "", errors.New("DIVOOM_WEB_PORT must be an integer between 1 and 65535")
	}
	return strconv.Itoa(number), nil
}

func LoadWebPassword() (string, error) {
	password := os.Getenv("DIVOOM_WEB_PASSWORD")
	if strings.TrimSpace(password) == "" {
		return "", errors.New("set DIVOOM_WEB_PASSWORD for web login")
	}
	if len(password) > 1024 {
		return "", errors.New("DIVOOM_WEB_PASSWORD must be at most 1024 bytes")
	}
	return password, nil
}

// LoadDataDirectory supplies the shared root for all persistent application data.
func LoadDataDirectory() string {
	if dir := os.Getenv("DATA_DIR"); dir != "" {
		return dir
	}
	return "var/lib/divoom_gallery_web/"
}

func LoadWebAddress() (string, error) {
	port, err := loadWebPort()
	if err != nil {
		return "", err
	}
	host := os.Getenv("DIVOOM_WEB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	return net.JoinHostPort(host, port), nil
}

// LoadMetrics disables reporting when no endpoint URL is configured.
func LoadMetrics() (Metrics, error) {
	m := Metrics{
		URL:      strings.TrimSpace(os.Getenv("DIVOOM_METRICS_URL")),
		Username: os.Getenv("DIVOOM_METRICS_USERNAME"),
		Password: os.Getenv("DIVOOM_METRICS_PASSWORD"),
	}
	if m.URL == "" {
		return m, nil
	}
	if strings.TrimSpace(m.Username) == "" || m.Password == "" || strings.ContainsAny(m.Username, ":\r\n") || strings.ContainsAny(m.Password, "\r\n") {
		return m, errors.New("set DIVOOM_METRICS_USERNAME and DIVOOM_METRICS_PASSWORD without line breaks; username must not contain a colon")
	}
	u, err := url.Parse(m.URL)
	if err != nil || u.Hostname() == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return m, errors.New("DIVOOM_METRICS_URL must be an HTTP(S) import URL without credentials, query, or fragment")
	}
	return m, nil
}
