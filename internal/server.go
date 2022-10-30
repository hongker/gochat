package internal

import (
	"context"
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/ebar-go/znet"
	"gochat/internal/bucket"
	"gochat/internal/http"
	"gochat/internal/socket"
	"log"
)

type Server struct{}

func (server *Server) Run(stopCh <-chan struct{}) (err error) {
	log.Println("server started")
	b := bucket.NewBucket()

	httpContext, httpCancel := context.WithCancel(context.Background())
	defer httpCancel()
	go func() {
		defer runtime.HandleCrash()
		ego.NewHTTPServer(":8080").
			EnableReleaseMode().
			EnableCorsMiddleware().
			EnablePprofHandler().
			EnableAvailableHealthCheck().
			RegisterRouteLoader(http.NewHandler(b).Install).
			Serve(httpContext.Done())
	}()

	handler := socket.NewHandler(b)
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
