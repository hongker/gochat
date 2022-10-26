package socket

import (
	"github.com/ebar-go/znet"
	"github.com/ebar-go/znet/codec"
	"gochat/api"
	"gochat/internal/application"
)

// listSession return a list of user sessions
func (handler *Handler) listSession(ctx *znet.Context, req *api.SessionListRequest) (resp *api.SessionListResponse, err error) {
	uid := handler.currentUser(ctx)
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

// sendMessage sends a message to receiver
func (handler *Handler) sendMessage(ctx *znet.Context, req *api.MessageSendRequest) (resp *api.MessageSendResponse, err error) {
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
		// save user session
		handler.sessionApp.SaveSession(ctx, uid, &application.Session{ID: req.Target, Title: "", Last: msg})
		handler.sessionApp.SaveSession(ctx, req.Target, &application.Session{ID: uid, Title: "", Last: msg})
	}
	return
}

// queryMessage query session history message
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
