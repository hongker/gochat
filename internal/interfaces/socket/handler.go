package socket

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/ebar-go/znet"
	"gochat/internal/application"
	"gochat/internal/domain/bucket"
	"gochat/internal/domain/dto"
	"gochat/pkg/cmap"
	"time"
)

type Handler struct {
	timers            *cmap.Container[string, *time.Timer]
	bucket            *bucket.Bucket
	heartbeatInterval time.Duration
	sessionApp        *application.SessionApplication
	messageApp        *application.MessageApplication
	channelApp        *application.ChannelApplication
	userApp           *application.UserApplication
	total             int
}

func NewHandler() *Handler {
	return &Handler{
		bucket:            bucket.NewBucket(),
		timers:            cmap.NewContainer[string, *time.Timer](),
		heartbeatInterval: time.Second * 60,
		sessionApp:        application.NewSessionApplication(),
		messageApp:        application.NewMessageApplication(),
		channelApp:        application.NewChannelApplication(),
		userApp:           application.NewUserApplication(),
	}
}

func (handler *Handler) OnConnect(conn *znet.Connection) {
	handler.total++
	component.Provider().Logger().Infof("[%s] connected:%s, total=%d", conn.ID(), conn.IP(), handler.total)

	// 给每个客户端设置一个定时器，如果这段时间内没有发送心跳，则自动关闭连接
	timer := time.NewTimer(handler.heartbeatInterval)
	go func() {
		defer runtime.HandleCrash()
		<-timer.C

		conn.Close()
	}()
	handler.timers.Set(conn.ID(), timer)

}
func (handler *Handler) OnDisconnect(conn *znet.Connection) {
	handler.total--
	component.Provider().Logger().Infof("[%s] Disconnected:%s", conn.ID(), conn.IP())

	// 清除定时器
	timer, _ := handler.timers.Get(conn.ID())
	if timer == nil {
		return
	}
	timer.Stop()
	handler.timers.Del(conn.ID())
}

func (handler *Handler) heartbeat(ctx *znet.Context, req *dto.HeartbeatRequest) (resp *dto.HeartbeatResponse, err error) {
	timer, _ := handler.timers.Get(ctx.Conn().ID())
	if timer == nil {
		return
	}

	timer.Reset(handler.heartbeatInterval)
	resp = &dto.HeartbeatResponse{ServerTime: time.Now().UnixMilli()}
	return
}
