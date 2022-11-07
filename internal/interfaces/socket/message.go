package socket

import (
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/znet"
	"github.com/ebar-go/znet/codec"
	"gochat/api"
	"gochat/internal/domain/constant"
	"gochat/internal/domain/dto"
	"gochat/internal/domain/types"
)

// listSession return a list of user sessions
func (handler *Handler) listSession(ctx *znet.Context, req *dto.SessionListRequest) (resp *dto.SessionListResponse, err error) {
	uid := handler.currentUser(ctx)
	// private user sessions
	sessions, err := handler.sessionApp.GetUserSessionList(ctx, uid)
	if err != nil {
		return
	}

	resp = &dto.SessionListResponse{Items: make([]dto.Session, 0, len(sessions))}
	for _, session := range sessions {
		item := dto.Session{
			ID:    session.ID,
			Title: session.Title,
			Type:  constant.SessionTypeUser,
		}
		if session.Last != nil {
			item.Last = dto.Message{
				ID:          session.Last.ID,
				Content:     session.Last.Content,
				ContentType: session.Last.ContentType,
				CreatedAt:   session.Last.CreatedAt,
			}
		}
		resp.Items = append(resp.Items, item)
	}

	// group sessions
	channels := handler.channelApp.GetJoined(ctx, handler.bucket, uid)
	for _, channel := range channels {
		item := dto.Session{
			ID:    channel.ID,
			Title: channel.Name,
			Type:  constant.SessionTypeGroup,
		}
		last := handler.messageApp.GetLast(ctx, item.ID)
		if last != nil {
			item.Last = dto.Message{
				ID:          last.ID,
				Content:     last.Content,
				ContentType: last.ContentType,
				CreatedAt:   last.CreatedAt,
			}
		}
		resp.Items = append(resp.Items, item)
	}

	return
}

func (handler *Handler) notify(packet codec.Codec, uid string, sender *types.User, msg *types.Message) (err error) {
	userSession := handler.bucket.GetSession(uid)
	if userSession == nil {
		err = errors.NotFound("user is not online")
		return
	}

	buf, err := packet.Pack(dto.Message{
		ID:           msg.ID,
		SessionID:    msg.SessionID,
		SessionTitle: sender.Name,
		Content:      msg.Content,
		ContentType:  msg.ContentType,
		CreatedAt:    msg.CreatedAt,
		Sender:       dto.User{ID: sender.ID, Name: sender.Name, Avatar: sender.Avatar},
	})
	if err != nil {
		return
	}

	userSession.Send(buf)
	return
}

// sendMessage sends a message to receiver
func (handler *Handler) sendMessage(ctx *znet.Context, req *dto.MessageSendRequest) (resp *dto.MessageSendResponse, err error) {
	uid := handler.currentUser(ctx)
	if req.Target == uid {
		return nil, errors.InvalidParam("can't send message to yourself")
	}
	var sender, receiver *types.User
	sender, err = handler.userApp.Get(ctx, uid)
	if err != nil {
		return
	}
	receiver, err = handler.userApp.Get(ctx, req.Target)
	if err != nil {
		return
	}

	// save message
	sessionId := handler.sessionApp.BuildUserSessionId(uid, req.Target)
	msg := &types.Message{
		SessionID:   sessionId,
		SessionType: constant.SessionTypeUser,
		Content:     req.Content,
		ContentType: req.ContentType,
		Target:      req.Target,
		Sender:      sender.ID,
	}
	handler.messageApp.Save(sessionId, msg)

	// send message
	packet := codec.Factory().NewWithHeader(codec.Header{Operate: api.OperatePushMessage, ContentType: ctx.Header().ContentType})
	_ = handler.notify(packet, req.Target, sender, msg)
	_ = handler.notify(packet, uid, sender, msg)

	// save user session
	_ = handler.sessionApp.SaveSession(ctx, uid, &types.Session{ID: sessionId, Title: receiver.Name, Last: msg})
	_ = handler.sessionApp.SaveSession(ctx, req.Target, &types.Session{ID: sessionId, Title: sender.Name, Last: msg})

	handler.userApp.SaveContact(ctx, uid, req.Target)
	handler.userApp.SaveContact(ctx, req.Target, uid)
	return
}

// queryMessage query session history message
func (handler *Handler) queryMessage(ctx *znet.Context, req *dto.MessageQueryRequest) (resp *dto.MessageQueryResponse, err error) {
	items, err := handler.messageApp.Query(ctx, req.SessionID)
	if err != nil {
		return
	}

	resp = &dto.MessageQueryResponse{SessionID: req.SessionID, Items: make([]dto.Message, len(items))}
	for idx, item := range items {
		message := dto.Message{
			ID:          item.ID,
			Content:     item.Content,
			CreatedAt:   item.CreatedAt,
			ContentType: item.ContentType,
		}
		sender, lastErr := handler.userApp.Get(ctx, item.Sender)
		if lastErr == nil {
			message.Sender = dto.User{ID: sender.ID, Name: sender.Name, Avatar: sender.Avatar}
		}
		resp.Items[idx] = message
	}
	return
}
