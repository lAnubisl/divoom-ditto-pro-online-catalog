package config

// DeviceAPI configures the device HTTP API independently of gallery credentials.
type DeviceAPI struct {
	BaseURL string
	APIKey  string
}
