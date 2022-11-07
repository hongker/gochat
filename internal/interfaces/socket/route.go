package socket

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/znet"
	"gochat/api"
	"gochat/internal/domain/dto"
)

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

	router.Route(api.OperateQueryChannel, znet.StandardHandler[dto.ChannelQueryRequest, dto.ChannelQueryResponse](handler.listChannel))
	router.Route(api.OperateJoinChannel, znet.StandardHandler[dto.ChannelJoinRequest, dto.ChannelJoinResponse](handler.joinChannel))
	router.Route(api.OperateLeaveChannel, znet.StandardHandler[dto.ChannelLeaveRequest, dto.ChannelLeaveResponse](handler.leaveChannel))
	router.Route(api.OperateCreateChannel, znet.StandardHandler[dto.ChannelCreateRequest, dto.ChannelCreateResponse](handler.createChannel))
	router.Route(api.OperateBroadcastChannel, znet.StandardHandler[dto.ChannelBroadcastRequest, dto.ChannelBroadcastResponse](handler.broadcastChannel))

	router.Route(api.OperateQueryMessage, znet.StandardHandler[dto.MessageQueryRequest, dto.MessageQueryResponse](handler.queryMessage))

	router.Route(api.OperateQueryContact, znet.StandardHandler[dto.ContactQueryRequest, dto.ContactQueryResponse](handler.getContacts))
}
