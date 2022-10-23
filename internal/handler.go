package internal

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/ebar-go/znet"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

const (
	OperateHeartbeat = 1
	OperateLogin     = 2
)

type Handler struct {
	rmw               sync.RWMutex
	timers            map[string]*time.Timer
	heartbeatInterval time.Duration
}

func (handler *Handler) Install(router *znet.Router) {
	router.Route(OperateLogin, znet.StandardHandler[LoginRequest, LoginResponse](handler.login))
	router.Route(OperateHeartbeat, znet.StandardHandler[HeartbeatRequest, HeartbeatResponse](handler.heartbeat))
}

func (handler *Handler) OnConnect(conn *znet.Connection) {
	timer := time.NewTimer(handler.heartbeatInterval)
	go func() {
		defer runtime.HandleCrash()
		<-timer.C

		conn.Close()
	}()

}
func (handler *Handler) OnDisconnect(conn *znet.Connection) {
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
	if ctx.Request().Header.Operate == OperateLogin {
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

func (handler *Handler) login(ctx *znet.Context, req *LoginRequest) (resp *LoginResponse, err error) {
	resp = &LoginResponse{UID: uuid.NewV4().String()}
	ctx.Conn().Property().Set("uid", resp.UID)
	ctx.Conn().Property().Set("username", req.Name)
	return
}

func (handler *Handler) heartbeat(ctx *znet.Context, req *HeartbeatRequest) (resp *HeartbeatResponse, err error) {
	handler.rmw.RLock()
	timer := handler.timers[ctx.Conn().ID()]
	handler.rmw.RUnlock()
	if timer == nil {
		return
	}

	timer.Reset(handler.heartbeatInterval)
	return
}

func NewHandler() *Handler {
	return &Handler{
		timers:            map[string]*time.Timer{},
		heartbeatInterval: time.Minute,
	}
}
