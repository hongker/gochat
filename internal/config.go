package internal

type Config struct{}

func (config *Config) New() *Server {
	return &Server{}
}

func DefaultConfig() *Config {
	return &Config{}
}
