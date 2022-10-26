package internal

import (
	"context"
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/ebar-go/znet"
	"gochat/internal/http"
	"gochat/internal/socket"
	"log"
)

type Server struct{}

func (server *Server) Run(stopCh <-chan struct{}) (err error) {
	log.Println("server started")

	httpContext, httpCancel := context.WithCancel(context.Background())
	defer httpCancel()
	go func() {
		defer runtime.HandleCrash()
		ego.NewHTTPServer(":8080").EnableReleaseMode().RegisterRouteLoader(http.NewHandler().Install).
			Serve(httpContext.Done())
	}()

	handler := socket.NewHandler()
	instance := znet.New(func(options *znet.Options) {
		options.OnConnect = handler.OnConnect
		options.OnDisconnect = handler.OnDisconnect
		options.Middlewares = append(options.Middlewares, handler.WriteRequestLog, handler.CheckLogin)
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
