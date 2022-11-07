package application

import (
	"context"
	"fmt"
	"gochat/internal/domain/types"
	"gochat/pkg/cmap"
)

type SessionApplication struct {
	collection *cmap.Container[string, []*types.Session]
}

func (app *SessionApplication) GetUserSessionList(ctx context.Context, uid string) ([]*types.Session, error) {
	items, _ := app.collection.Get(uid)
	return items, nil
}

func (app *SessionApplication) BuildUserSessionId(uid, targetId string) string {
	if uid > targetId {
		return fmt.Sprintf("%s:%s", uid, targetId)
	}
	return fmt.Sprintf("%s:%s", targetId, uid)
}

func (app *SessionApplication) SaveSession(ctx context.Context, uid string, session *types.Session) (err error) {
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
		collection: cmap.NewContainer[string, []*types.Session](),
	}
}
