package application

import (
	"context"
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
