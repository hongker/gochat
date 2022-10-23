package application

import (
	"context"
	"github.com/ebar-go/znet/codec"
	uuid "github.com/satori/go.uuid"
	"gochat/internal/bucket"
	"sync"
	"time"
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

	app.bucket.SubscribeChannel(app.bucket.AddChannel(id), app.bucket.GetSession(uid))

	return
}

func (app *ChannelApplication) Join(ctx context.Context, id string, uid string) (err error) {
	channel := app.bucket.GetChannel(id)
	if channel == nil {
		return
	}

	app.bucket.SubscribeChannel(channel, app.bucket.GetSession(uid))

	return
}

func (app *ChannelApplication) Leave(ctx context.Context, id string, uid string) (err error) {
	channel := app.bucket.GetChannel(id)
	if channel == nil {
		return
	}

	app.bucket.UnsubscribeChannel(channel, app.bucket.GetSession(uid))
	return
}

func (app *ChannelApplication) Broadcast(ctx context.Context, msg *Message, codec codec.Codec, packet *codec.Packet) (err error) {
	channel := app.bucket.GetChannel(msg.Target)
	if channel == nil {
		return
	}

	msg.ID = uuid.NewV4().String()
	msg.CreatedAt = time.Now().UnixMilli()

	buf, err := codec.Pack(packet, msg)
	if err != nil {
		return
	}

	app.bucket.BroadcastChannel(channel, buf)
	return
}

func NewChannelApplication(bucket *bucket.Bucket) *ChannelApplication {
	return &ChannelApplication{
		bucket:   bucket,
		channels: map[string]*Channel{},
	}
}
