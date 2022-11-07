package application

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"gochat/internal/domain/types"
	"gochat/pkg/cmap"
	"time"
)

type MessageApplication struct {
	messages               *cmap.Container[string, []*types.Message]
	maxSessionMessageCount int
}

func (app *MessageApplication) Query(ctx context.Context, sessionID string) (items []*types.Message, err error) {
	items, _ = app.messages.Get(sessionID)
	return
}

func (app *MessageApplication) GetLast(ctx context.Context, sessionID string) (item *types.Message) {
	items, _ := app.messages.Get(sessionID)
	if len(items) > 0 {
		item = items[len(items)-1]
	}
	return
}

func (app *MessageApplication) Save(sessionId string, msg *types.Message) {
	if msg.ID == "" {
		msg.ID = uuid.NewV4().String()
		msg.CreatedAt = time.Now().UnixMilli()
	}
	items, _ := app.messages.Get(sessionId)
	if len(items) == 0 {
		items = make([]*types.Message, 0, app.maxSessionMessageCount)
	}
	items = append(items, msg)
	if total := len(items); total > app.maxSessionMessageCount {
		items = items[total-app.maxSessionMessageCount : total]
	}
	app.messages.Set(sessionId, items)

}

func NewMessageApplication() *MessageApplication {
	return &MessageApplication{
		messages:               cmap.NewContainer[string, []*types.Message](),
		maxSessionMessageCount: 100,
	}
}
