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

	handler := NewHandler()
	instance := znet.New(func(options *znet.Options) {
		options.OnConnect = handler.OnConnect
		options.OnDisconnect = handler.OnDisconnect
		options.Middlewares = append(options.Middlewares, handler.DebugLog, handler.CheckLogin)
	})

	handler.Install(instance.Router())

	instance.ListenTCP(":8081")
	instance.ListenWebsocket(":8082")
	if err = instance.Run(stopCh); err != nil {
		return
	}

	server.shutdown()
	return

}

func (server *Server) shutdown() {
	log.Println("server stopped")
}
