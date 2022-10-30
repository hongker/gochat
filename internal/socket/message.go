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
	sessions, err := handler.sessionApp.GetSessionList(ctx, uid)
	if err != nil {
		return
	}

	resp = &dto.SessionListResponse{Items: make([]dto.Session, 0, len(sessions))}
	for _, session := range sessions {
		item := dto.Session{
			ID:    session.ID,
			Title: session.Title,
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
	return
}

// sendMessage sends a message to receiver
func (handler *Handler) sendMessage(ctx *znet.Context, req *dto.MessageSendRequest) (resp *dto.MessageSendResponse, err error) {
	uid := handler.currentUser(ctx)
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
		var sender, receiver *application.User
		sender, err = handler.userApp.Get(ctx, uid)
		if err != nil {
			return
		}
		receiver, err = handler.userApp.Get(ctx, req.Target)
		if err != nil {
			return
		}
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
		resp.Items[idx] = dto.Message{
			ID:          item.ID,
			Content:     item.Content,
			CreatedAt:   item.CreatedAt,
			ContentType: item.ContentType,
		}
	}
	return
}
