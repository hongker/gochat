package application

import (
	"context"
	"github.com/ebar-go/ego/errors"
	uuid "github.com/satori/go.uuid"
	"gochat/internal/domain/bucket"
	"gochat/pkg/cmap"
)

type ChannelApplication struct {
	collection *cmap.Container[string, *Channel]
}

type Channel struct {
	ID    string
	Name  string
	Owner string
}

func (app *ChannelApplication) List(ctx context.Context) []*Channel {
	items := make([]*Channel, 0)
	app.collection.Iterator(func(key string, item *Channel) {
		items = append(items, item)
	})
	return items
}

func (app *ChannelApplication) GetJoined(ctx context.Context, bucket *bucket.Bucket, uid string) []*Channel {
	session := bucket.GetSession(uid)
	if session == nil {
		return nil
	}

	items := make([]*Channel, 0)
	for _, id := range session.Channels {
		item, exist := app.collection.Get(id)
		if !exist {
			continue
		}
		items = append(items, item)
	}
	return items
}

func (app *ChannelApplication) Create(ctx context.Context, uid, name string) (channel *Channel, err error) {
	id := uuid.NewV4().String()
	channel = &Channel{ID: id, Name: name, Owner: uid}
	app.collection.Set(id, channel)
	return
}

func (app *ChannelApplication) Get(ctx context.Context, id string) (*Channel, error) {
	channel, exist := app.collection.Get(id)
	if !exist {
		return nil, errors.NotFound("channel not found")
	}
	return channel, nil
}

func NewChannelApplication() *ChannelApplication {
	return &ChannelApplication{
		collection: cmap.NewContainer[string, *Channel](),
	}
}
