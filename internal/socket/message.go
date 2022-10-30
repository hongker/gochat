package socket

import (
	"github.com/ebar-go/znet"
	"github.com/ebar-go/znet/codec"
	"gochat/api"
	"gochat/internal/application"
	"gochat/internal/domain/dto"
)

// listSession return a list of user sessions
func (handler *Handler) listSession(ctx *znet.Context, req *dto.SessionListRequest) (resp *dto.SessionListResponse, err error) {
	uid := handler.currentUser(ctx)
	// private user sessions
	sessions, err := handler.sessionApp.GetSessionList(ctx, uid)
	if err != nil {
		return
	}

	resp = &dto.SessionListResponse{Items: make([]dto.Session, 0, len(sessions))}
	for _, session := range sessions {
		item := dto.Session{
			ID:    session.ID,
			Title: session.Title,
			Type:  api.SessionTypeUser,
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
	channels := handler.channelApp.GetJoined(ctx, uid)
	for _, channel := range channels {
		item := dto.Session{
			ID:    channel.ID,
			Title: channel.Name,
			Type:  api.SessionTypeGroup,
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

// sendMessage sends a message to receiver
func (handler *Handler) sendMessage(ctx *znet.Context, req *dto.MessageSendRequest) (resp *dto.MessageSendResponse, err error) {
	uid := handler.currentUser(ctx)
	var sender, receiver *application.User
	sender, err = handler.userApp.Get(ctx, uid)
	if err != nil {
		return
	}
	receiver, err = handler.userApp.Get(ctx, req.Target)
	if err != nil {
		return
	}

	packet := &codec.Packet{Header: codec.Header{Operate: api.OperatePushMessage, ContentType: ctx.Request().Header.ContentType}}
	msg, err := handler.messageApp.Send(ctx, sender, receiver, req, codec.Default(), packet)
	if err == nil {
		// save user session
		handler.sessionApp.SaveSession(ctx, uid, &application.Session{ID: req.Target, Title: receiver.Name, Last: msg})
		handler.sessionApp.SaveSession(ctx, req.Target, &application.Session{ID: uid, Title: sender.Name, Last: msg})

		handler.userApp.SaveContact(ctx, uid, req.Target)
		handler.userApp.SaveContact(ctx, req.Target, uid)
	}
	return
}

// queryMessage query session history message
func (handler *Handler) queryMessage(ctx *znet.Context, req *dto.MessageQueryRequest) (resp *dto.MessageQueryResponse, err error) {
	items, err := handler.messageApp.Query(ctx, req.SessionID)
	if err != nil {
		return
	}

	resp = &dto.MessageQueryResponse{Items: make([]dto.Message, len(items))}
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
