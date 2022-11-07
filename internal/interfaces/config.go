package interfaces

import (
	"gochat/internal/interfaces/socket"
	"time"
)

type Config struct {
	PprofEnable       bool
	HealthCheckEnable bool
	SocketOptions     socket.Options
}

func (config *Config) New() *Server {
	return &Server{
		config: config,
	}
}

func DefaultConfig() *Config {
	return &Config{
		PprofEnable: true,
		SocketOptions: socket.Options{
			BucketQueueCount:  32,
			BucketQueueSize:   100,
			HeartbeatInterval: time.Minute,
		},
	}
}
