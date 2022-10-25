package socket

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/ebar-go/znet"
	"github.com/ebar-go/znet/codec"
	"gochat/api"
	"gochat/internal/application"
	"gochat/internal/bucket"
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

func (handler *Handler) Install(router *znet.Router) {
	router.OnError(func(ctx *znet.Context, err error) {
		component.Provider().Logger().Errorf("[%s] error: %v", ctx.Conn().ID(), err)
	})
	router.Route(api.OperateLogin, znet.StandardHandler[api.LoginRequest, api.LoginResponse](handler.login))
	router.Route(api.OperateHeartbeat, znet.StandardHandler[api.HeartbeatRequest, api.HeartbeatResponse](handler.heartbeat))
	router.Route(api.OperateListSession, znet.StandardHandler[api.SessionListRequest, api.SessionListResponse](handler.listSession))
	router.Route(api.OperateSendMessage, znet.StandardHandler[api.MessageSendRequest, api.MessageSendResponse](handler.sendMessage))

	router.Route(api.OperateJoinChannel, znet.StandardHandler[api.ChannelJoinRequest, api.ChannelJoinResponse](handler.joinChannel))
	router.Route(api.OperateLeaveChannel, znet.StandardHandler[api.ChannelLeaveRequest, api.ChannelLeaveResponse](handler.leaveChannel))
	router.Route(api.OperateCreateChannel, znet.StandardHandler[api.ChannelCreateRequest, api.ChannelCreateResponse](handler.createChannel))
	router.Route(api.OperateBroadcastChannel, znet.StandardHandler[api.ChannelBroadcastRequest, api.ChannelBroadcastResponse](handler.broadcastChannel))

	router.Route(api.OperateQueryMessage, znet.StandardHandler[api.MessageQueryRequest, api.MessageQueryResponse](handler.queryMessage))
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

func (handler *Handler) DebugLog(ctx *znet.Context) {
	component.Provider().Logger().Infof("[%s] request log: header=%+v, body=%v ", ctx.Conn().ID(), ctx.Request().Header, string(ctx.Request().Body))
	ctx.Next()
}

func (handler *Handler) CheckLogin(ctx *znet.Context) {
	if ctx.Request().Header.Operate == api.OperateLogin {
		ctx.Next()
		return
	}

	_, err := handler.getCurrentUser(ctx)
	if err != nil {
		component.Provider().Logger().Error(err.Error())
		ctx.Abort()
		return
	}
	ctx.Next()
}

func (handler *Handler) getCurrentUser(ctx *znet.Context) (string, error) {
	uid, exist := ctx.Conn().Property().Get("uid")
	if !exist {
		return "", errors.Unauthorized("unauthorized")
	}
	return uid.(string), nil
}

func (handler *Handler) login(ctx *znet.Context, req *api.LoginRequest) (resp *api.LoginResponse, err error) {
	user := &application.User{Name: req.Name}
	err = handler.userApp.Login(ctx, user)
	if err != nil {
		return
	}
	resp = &api.LoginResponse{UID: user.ID}
	ctx.Conn().Property().Set("uid", user.ID)
	handler.bucket.AddSession(bucket.NewSession(user.ID, ctx.Conn()))
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

func (handler *Handler) listSession(ctx *znet.Context, req *api.SessionListRequest) (resp *api.SessionListResponse, err error) {
	uid, err := handler.getCurrentUser(ctx)
	if err != nil {
		return
	}
	sessions, err := handler.sessionApp.GetSessionList(ctx, uid)
	if err != nil {
		return
	}

	resp = &api.SessionListResponse{Items: make([]api.Session, 0, len(sessions))}
	for _, session := range sessions {
		resp.Items = append(resp.Items, api.Session{ID: session.ID, Title: session.Title})
	}
	return
}

func (handler *Handler) sendMessage(ctx *znet.Context, req *api.MessageSendRequest) (resp *api.MessageSendResponse, err error) {
	uid, err := handler.getCurrentUser(ctx)
	if err != nil {
		return
	}
	msg := &application.Message{
		Content:     req.Content,
		ContentType: req.ContentType,
		Target:      req.Target,
		Sender:      uid,
	}

	packet := &codec.Packet{Header: codec.Header{Operate: api.OperatePushMessage, ContentType: ctx.Request().Header.ContentType}}
	err = handler.messageApp.Send(ctx, msg, codec.Default(), packet)
	if err == nil {
		// save user session
		handler.sessionApp.SaveSession(ctx, uid, &application.Session{ID: req.Target, Title: "", Last: msg})
		handler.sessionApp.SaveSession(ctx, req.Target, &application.Session{ID: uid, Title: "", Last: msg})
	}
	return
}

func (handler *Handler) createChannel(ctx *znet.Context, req *api.ChannelCreateRequest) (resp *api.ChannelCreateResponse, err error) {

	uid, _ := handler.getCurrentUser(ctx)
	channel, err := handler.channelApp.Create(ctx, uid, req.Name)
	if err != nil {
		return
	}
	resp = &api.ChannelCreateResponse{ID: channel.ID}
	return
}

func (handler *Handler) joinChannel(ctx *znet.Context, req *api.ChannelJoinRequest) (resp *api.ChannelJoinResponse, err error) {
	uid, _ := handler.getCurrentUser(ctx)
	err = handler.channelApp.Join(ctx, req.ID, uid)
	return
}
func (handler *Handler) leaveChannel(ctx *znet.Context, req *api.ChannelLeaveRequest) (resp *api.ChannelLeaveResponse, err error) {
	uid, _ := handler.getCurrentUser(ctx)
	err = handler.channelApp.Leave(ctx, req.ID, uid)
	return
}

func (handler *Handler) broadcastChannel(ctx *znet.Context, req *api.ChannelBroadcastRequest) (resp *api.ChannelBroadcastResponse, err error) {
	packet := &codec.Packet{Header: codec.Header{Operate: api.OperatePushMessage, ContentType: ctx.Request().Header.ContentType}}

	uid, _ := handler.getCurrentUser(ctx)

	msg := &application.Message{
		Content:     req.Content,
		ContentType: req.ContentType,
		Target:      req.Target,
		Sender:      uid,
	}
	err = handler.channelApp.Broadcast(ctx, msg, codec.Default(), packet)

	if err != nil {
		return
	}

	handler.messageApp.Save(req.Target, msg)
	return
}

func (handler *Handler) queryMessage(ctx *znet.Context, req *api.MessageQueryRequest) (resp *api.MessageQueryResponse, err error) {
	items, err := handler.messageApp.Query(ctx, req.SessionID)
	if err != nil {
		return
	}

	resp = &api.MessageQueryResponse{Items: make([]api.Message, len(items))}
	for idx, item := range items {
		resp.Items[idx] = api.Message{
			ID:          item.ID,
			Content:     item.Content,
			CreatedAt:   item.CreatedAt,
			ContentType: item.ContentType,
		}
	}
	return
}

func NewHandler() *Handler {
	b := bucket.NewBucket()
	return &Handler{
		bucket:            b,
		timers:            map[string]*time.Timer{},
		heartbeatInterval: time.Minute,
		sessionApp:        application.NewSessionApplication(),
		messageApp:        application.NewMessageApplication(b),
		channelApp:        application.NewChannelApplication(b),
		userApp:           application.NewUserApplication(),
	}
}