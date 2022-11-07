package interfaces

type Config struct {
	PprofEnable       bool
	HealthCheckEnable bool
}

func (config *Config) New() *Server {
	return &Server{
		config: config,
	}
}

func DefaultConfig() *Config {
	return &Config{
		PprofEnable: true,
	}
}
