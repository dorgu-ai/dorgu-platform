package server

// Config holds the server configuration.
type Config struct {
	Port       int
	KubeConfig string
	Context    string
	Host       string // default "localhost"
}

// NewConfig returns a default configuration.
func NewConfig() *Config {
	return &Config{
		Port: 8080,
		Host: "localhost",
	}
}
