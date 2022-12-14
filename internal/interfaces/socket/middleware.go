package socket

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/znet"
	"gochat/api"
)

// WriteRequestLog represents print request log middleware
func (handler *Handler) WriteRequestLog(ctx *znet.Context) {
	component.Provider().Logger().Infof("[%s] request log: header=%+v, body=%v ", ctx.Conn().ID(), ctx.Header(), string(ctx.Request()))
	ctx.Next()
}

// CheckLogin validates the login credentials
func (handler *Handler) CheckLogin(ctx *znet.Context) {
	// skip when operate is Login
	if ctx.Header().Operate == api.OperateConnect {
		ctx.Next()
		return
	}

	// validate login credentials
	if uid := handler.currentUser(ctx); uid == "" {
		component.Provider().Logger().Error("unauthorized")
		ctx.Abort()
		return
	}
	ctx.Next()
}
