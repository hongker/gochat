package internal

import "github.com/ebar-go/znet"

const (
	OperateLogin     = 1
	OperateHeartbeat = 2
)

type Handler struct{}

func (handler *Handler) Install(router *znet.Router) {
	router.Route(OperateLogin, znet.StandardHandler[LoginRequest, LoginResponse](handler.login))
	router.Route(OperateHeartbeat, znet.StandardHandler[HeartbeatRequest, HeartbeatResponse](handler.heartbeat))
}

func (handler *Handler) login(ctx *znet.Context, req *LoginRequest) (resp *LoginResponse, err error) {
	return
}

func (handler *Handler) heartbeat(ctx *znet.Context, req *HeartbeatRequest) (resp *HeartbeatResponse, err error) {
	return
}
