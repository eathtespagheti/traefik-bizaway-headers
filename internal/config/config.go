package config

// Config the plugin configuration.
type Config struct {
	Headers HeaderConfig `json:"headers,omitempty"`
}

// HeaderConfig defines the structure for request and response headers.
type HeaderConfig struct {
	Request  map[string]string `json:"request,omitempty"`
	Response map[string]string `json:"response,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers: HeaderConfig{
			Request:  make(map[string]string),
			Response: make(map[string]string),
		},
	}
}
