package application

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"gochat/internal/bucket"
	"sync"
)

type ChannelApplication struct {
	bucket   *bucket.Bucket
	rmu      sync.RWMutex
	channels map[string]*Channel
}

type Channel struct {
	ID    string
	Name  string
	Owner string
}

func (app *ChannelApplication) Create(ctx context.Context, uid, name string) (channel *Channel, err error) {
	id := uuid.NewV4().String()
	channel = &Channel{ID: id, Name: name, Owner: uid}
	app.rmu.Lock()
	app.channels[id] = channel
	app.rmu.Unlock()
	app.bucket.AddChannel(id)

	return
}

func (app *ChannelApplication) Get(ctx context.Context, id string) (channel *Channel, err error) {

	return
}

func NewChannelApplication(bucket *bucket.Bucket) *ChannelApplication {
	return &ChannelApplication{
		bucket:   bucket,
		channels: map[string]*Channel{},
	}
}
