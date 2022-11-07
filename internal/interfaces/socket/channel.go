package socket

import (
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/znet"
	"github.com/ebar-go/znet/codec"
	"gochat/api"
	"gochat/internal/domain/dto"
	"gochat/internal/domain/types"
)

func (handler *Handler) listChannel(ctx *znet.Context, req *dto.ChannelQueryRequest) (resp *dto.ChannelQueryResponse, err error) {
	items := handler.channelApp.List(ctx)
	resp = &dto.ChannelQueryResponse{Items: make([]dto.Channel, len(items))}
	for i, item := range items {
		resp.Items[i] = dto.Channel{
			ID:    item.ID,
			Name:  item.Name,
			Owner: item.Owner,
		}
	}
	return
}

func (handler *Handler) createChannel(ctx *znet.Context, req *dto.ChannelCreateRequest) (resp *dto.ChannelCreateResponse, err error) {
	uid := handler.currentUser(ctx)
	channel, err := handler.channelApp.Create(ctx, uid, req.Name)
	if err != nil {
		return
	}

	c := handler.bucket.AddChannel(channel.ID)
	handler.bucket.SubscribeChannel(c, handler.bucket.GetSession(uid))

	resp = &dto.ChannelCreateResponse{ID: channel.ID}
	for _, id := range req.Members {
		session := handler.bucket.GetSession(id)
		if session != nil {
			handler.bucket.SubscribeChannel(c, session)
		}
	}
	return
}

func (handler *Handler) joinChannel(ctx *znet.Context, req *dto.ChannelJoinRequest) (resp *dto.ChannelJoinResponse, err error) {
	uid := handler.currentUser(ctx)
	channel := handler.bucket.GetChannel(req.ID)
	if channel == nil {
		err = errors.NotFound("channel not found")
		return
	}
	handler.bucket.SubscribeChannel(channel, handler.bucket.GetSession(uid))
	return
}

func (handler *Handler) leaveChannel(ctx *znet.Context, req *dto.ChannelLeaveRequest) (resp *dto.ChannelLeaveResponse, err error) {
	uid := handler.currentUser(ctx)
	channel := handler.bucket.GetChannel(req.ID)
	if channel == nil {
		err = errors.NotFound("channel not found")
		return
	}

	handler.bucket.UnsubscribeChannel(channel, handler.bucket.GetSession(uid))
	return
}

func (handler *Handler) broadcastChannel(ctx *znet.Context, req *dto.ChannelBroadcastRequest) (resp *dto.ChannelBroadcastResponse, err error) {
	packet := codec.Factory().NewWithHeader(codec.Header{Operate: api.OperatePushMessage, ContentType: ctx.Header().ContentType})

	uid := handler.currentUser(ctx)
	sender, err := handler.userApp.Get(ctx, uid)
	if err != nil {
		return
	}

	msg := &types.Message{
		Content:     req.Content,
		ContentType: req.ContentType,
		Target:      req.Target,
		Sender:      uid,
	}

	handler.messageApp.Save(req.Target, msg)

	channel, err := handler.channelApp.Get(ctx, req.Target)
	if err != nil {
		return
	}

	buf, err := packet.Pack(dto.Message{
		ID:           msg.ID,
		SessionID:    msg.Target,
		SessionTitle: channel.Name,
		Content:      msg.Content,
		ContentType:  msg.ContentType,
		CreatedAt:    msg.CreatedAt,
		Sender:       dto.User{ID: sender.ID, Name: sender.Name, Avatar: sender.Avatar},
	})
	if err != nil {
		return
	}

	handler.bucket.BroadcastChannel(handler.bucket.GetChannel(req.Target), buf)

	return
}
