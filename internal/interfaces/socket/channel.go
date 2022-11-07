package socket

import (
	"github.com/ebar-go/znet"
	"github.com/ebar-go/znet/codec"
	uuid "github.com/satori/go.uuid"
	"gochat/api"
	"gochat/internal/application"
	"gochat/internal/domain/dto"
	"time"
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
	handler.channelApp.Join(ctx, channel.ID, req.Members...)
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
	packet := codec.Factory().NewWithHeader(codec.Header{Operate: api.OperatePushMessage, ContentType: ctx.Header().ContentType})

	uid := handler.currentUser(ctx)
	sender, err := handler.userApp.Get(ctx, uid)
	if err != nil {
		return
	}

	msg := &application.Message{
		ID:          uuid.NewV4().String(),
		Content:     req.Content,
		ContentType: req.ContentType,
		Target:      req.Target,
		Sender:      uid,
		CreatedAt:   time.Now().UnixMilli(),
	}

	handler.messageApp.Save(req.Target, msg)

	err = handler.channelApp.Broadcast(ctx, dto.Message{
		ID:          msg.ID,
		SessionID:   msg.Target,
		Content:     msg.Content,
		ContentType: msg.ContentType,
		CreatedAt:   msg.CreatedAt,
		Sender:      dto.User{ID: sender.ID, Name: sender.Name, Avatar: sender.Avatar},
	}, packet)

	return
}
