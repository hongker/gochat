package socket

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/ebar-go/znet"
	"gochat/api"
	"gochat/internal/application"
	"gochat/internal/bucket"
	"gochat/internal/domain/dto"
	"sync"
	"time"
)

type Handler struct {
	rmw               sync.RWMutex
	timers            map[string]*time.Timer
	bucket            *bucket.Bucket
	heartbeatInterval time.Duration
	sessionApp        *application.SessionApplication
	messageApp        *application.MessageApplication
	channelApp        *application.ChannelApplication
	userApp           *application.UserApplication
}

func NewHandler() *Handler {
	b := bucket.NewBucket()
	return &Handler{
		bucket:            b,
		timers:            map[string]*time.Timer{},
		heartbeatInterval: time.Minute * 10,
		sessionApp:        application.NewSessionApplication(),
		messageApp:        application.NewMessageApplication(b),
		channelApp:        application.NewChannelApplication(b),
		userApp:           application.NewUserApplication(),
	}
}

func (handler *Handler) Install(router *znet.Router) {
	router.OnError(func(ctx *znet.Context, err error) {
		component.Provider().Logger().Errorf("[%s] error: %v", ctx.Conn().ID(), err)
	})
	router.Route(api.OperateConnect, znet.StandardHandler[dto.ConnectRequest, dto.ConnectResponse](handler.connect))
	router.Route(api.OperateFindUserProfile, znet.StandardHandler[dto.IDRequest, dto.UserResponse](handler.findProfile))
	router.Route(api.OperateUpdateUserProfile, znet.StandardHandler[dto.UserUpdateRequest, dto.UserUpdateResponse](handler.updateProfile))
	router.Route(api.OperateHeartbeat, znet.StandardHandler[dto.HeartbeatRequest, dto.HeartbeatResponse](handler.heartbeat))
	router.Route(api.OperateListSession, znet.StandardHandler[dto.SessionListRequest, dto.SessionListResponse](handler.listSession))
	router.Route(api.OperateSendMessage, znet.StandardHandler[dto.MessageSendRequest, dto.MessageSendResponse](handler.sendMessage))

	router.Route(api.OperateJoinChannel, znet.StandardHandler[dto.ChannelJoinRequest, dto.ChannelJoinResponse](handler.joinChannel))
	router.Route(api.OperateLeaveChannel, znet.StandardHandler[dto.ChannelLeaveRequest, dto.ChannelLeaveResponse](handler.leaveChannel))
	router.Route(api.OperateCreateChannel, znet.StandardHandler[dto.ChannelCreateRequest, dto.ChannelCreateResponse](handler.createChannel))
	router.Route(api.OperateBroadcastChannel, znet.StandardHandler[dto.ChannelBroadcastRequest, dto.ChannelBroadcastResponse](handler.broadcastChannel))

	router.Route(api.OperateQueryMessage, znet.StandardHandler[dto.MessageQueryRequest, dto.MessageQueryResponse](handler.queryMessage))
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

	uid, exist := conn.Property().Get("uid")
	if !exist {
		return
	}
	session := handler.bucket.GetSession(uid.(string))
	for _, channel := range session.Channels {
		handler.bucket.UnsubscribeChannel(channel, session)
	}
}

func (handler *Handler) heartbeat(ctx *znet.Context, req *dto.HeartbeatRequest) (resp *dto.HeartbeatResponse, err error) {
	handler.rmw.RLock()
	timer := handler.timers[ctx.Conn().ID()]
	handler.rmw.RUnlock()
	if timer == nil {
		return
	}

	timer.Reset(handler.heartbeatInterval)
	resp = &dto.HeartbeatResponse{ServerTime: time.Now().UnixMilli()}
	return
}
