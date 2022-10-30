package socket

import (
	"github.com/ebar-go/znet"
	"github.com/ebar-go/znet/codec"
	"gochat/api"
	"gochat/internal/application"
	"gochat/internal/domain/dto"
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
	resp = &dto.ChannelCreateResponse{ID: channel.ID}
	return
}

func (handler *Handler) joinChannel(ctx *znet.Context, req *dto.ChannelJoinRequest) (resp *dto.ChannelJoinResponse, err error) {
	uid := handler.currentUser(ctx)
	err = handler.channelApp.Join(ctx, req.ID, uid)
	return
}
func (handler *Handler) leaveChannel(ctx *znet.Context, req *dto.ChannelLeaveRequest) (resp *dto.ChannelLeaveResponse, err error) {
	uid := handler.currentUser(ctx)
	err = handler.channelApp.Leave(ctx, req.ID, uid)
	return
}

func (handler *Handler) broadcastChannel(ctx *znet.Context, req *dto.ChannelBroadcastRequest) (resp *dto.ChannelBroadcastResponse, err error) {
	packet := &codec.Packet{Header: codec.Header{Operate: api.OperatePushMessage, ContentType: ctx.Request().Header.ContentType}}

	uid := handler.currentUser(ctx)

	msg := &application.Message{
		Content:     req.Content,
		ContentType: req.ContentType,
		Target:      req.Target,
		Sender:      uid,
	}

	handler.messageApp.Save(req.Target, msg)

	err = handler.channelApp.Broadcast(ctx, msg, codec.Default(), packet)

	return
}
