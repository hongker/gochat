package application

import (
	"context"
	"fmt"
	"gochat/pkg/cmap"
)

type SessionApplication struct {
	collection *cmap.Container[string, []*Session]
}

type Session struct {
	ID    string
	Title string
	Last  *Message
}

func (app *SessionApplication) GetSessionList(ctx context.Context, uid string) ([]*Session, error) {
	items, _ := app.collection.Get(uid)
	return items, nil
}

func (app *SessionApplication) BuildSessionId(uid, targetId string) string {
	if uid > targetId {
		return fmt.Sprintf("%s:%s", uid, targetId)
	}
	return fmt.Sprintf("%s:%s", targetId, uid)
}

func (app *SessionApplication) SaveSession(ctx context.Context, uid string, session *Session) (err error) {
	items, _ := app.collection.Get(uid)

	var exist bool
	for _, item := range items {
		if item.ID == session.ID {
			item.Title = session.Title
			item.Last = session.Last
			exist = true
			break
		}
	}

	if !exist {
		items = append(items, session)
	}
	app.collection.Set(uid, items)

	return
}

func NewSessionApplication() *SessionApplication {
	return &SessionApplication{
		collection: cmap.NewContainer[string, []*Session](),
	}
}
