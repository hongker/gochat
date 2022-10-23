package internal

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/ebar-go/znet"
	uuid "github.com/satori/go.uuid"
	"gochat/api"
	"sync"
	"time"
)

type Handler struct {
	rmw               sync.RWMutex
	timers            map[string]*time.Timer
	heartbeatInterval time.Duration
}

func (handler *Handler) Install(router *znet.Router) {
	router.Route(api.OperateLogin, znet.StandardHandler[api.LoginRequest, api.LoginResponse](handler.login))
	router.Route(api.OperateHeartbeat, znet.StandardHandler[api.HeartbeatRequest, api.HeartbeatResponse](handler.heartbeat))
}

func (handler *Handler) OnConnect(conn *znet.Connection) {
	component.Provider().Logger().Infof("[%s] connected", conn.ID())
	timer := time.NewTimer(handler.heartbeatInterval)
	go func() {
		defer runtime.HandleCrash()
		<-timer.C

		conn.Close()
	}()
	handler.rmw.Lock()
	handler.timers[conn.ID()] = timer
	handler.rmw.Unlock()

}
func (handler *Handler) OnDisconnect(conn *znet.Connection) {
	component.Provider().Logger().Infof("[%s] Disconnected", conn.ID())
	handler.rmw.RLock()
	timer := handler.timers[conn.ID()]
	handler.rmw.RUnlock()
	if timer == nil {
		return
	}
	timer.Stop()
	delete(handler.timers, conn.ID())
}

func (handler *Handler) CheckLogin(ctx *znet.Context) {
	if ctx.Request().Header.Operate == api.OperateLogin {
		ctx.Next()
		return
	}

	_, exist := ctx.Conn().Property().Get("uid")
	if !exist {
		component.Provider().Logger().Error("unauthorized")
		ctx.Abort()
		return
	}
	ctx.Next()
}

func (handler *Handler) login(ctx *znet.Context, req *api.LoginRequest) (resp *api.LoginResponse, err error) {
	resp = &api.LoginResponse{UID: uuid.NewV4().String()}
	ctx.Conn().Property().Set("uid", resp.UID)
	ctx.Conn().Property().Set("username", req.Name)
	return
}

func (handler *Handler) heartbeat(ctx *znet.Context, req *api.HeartbeatRequest) (resp *api.HeartbeatResponse, err error) {
	handler.rmw.RLock()
	timer := handler.timers[ctx.Conn().ID()]
	handler.rmw.RUnlock()
	if timer == nil {
		return
	}

	timer.Reset(handler.heartbeatInterval)
	resp = &api.HeartbeatResponse{ServerTime: time.Now().UnixMilli()}
	return
}

func NewHandler() *Handler {
	return &Handler{
		timers:            map[string]*time.Timer{},
		heartbeatInterval: time.Minute,
	}
}
