package internal

import (
	"github.com/ebar-go/znet"
	"log"
)

type Server struct{}

func (server *Server) Run(stopCh <-chan struct{}) (err error) {
	log.Println("server started")

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	instance := znet.New()

	instance.ListenTCP(":8081")
	if err = instance.Run(stopCh); err != nil {
		return
	}

	server.shutdown()
	return

}

func (server *Server) shutdown() {
	log.Println("server stopped")
}
