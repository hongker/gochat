package application

import (
	"context"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/znet/codec"
	uuid "github.com/satori/go.uuid"
	"gochat/internal/bucket"
	"sync"
	"time"
)

type MessageApplication struct {
	bucket   *bucket.Bucket
	rmw      sync.RWMutex
	messages map[string][]*Message
}

type Message struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	Target      string `json:"target"`
	Sender      string `json:"sender"`
	CreatedAt   int64  `json:"created_at"`
}

func (app *MessageApplication) Query(ctx context.Context, sessionID string) (items []*Message, err error) {
	app.rmw.RLock()
	items = app.messages[sessionID]
	app.rmw.RUnlock()
	return

}

func (app *MessageApplication) Send(ctx context.Context, msg *Message, codec codec.Codec, packet *codec.Packet) (err error) {
	receiverSession := app.bucket.GetSession(msg.Target)
	if receiverSession == nil {
		err = errors.NotFound("receiver not found")
		return
	}
	senderSession := app.bucket.GetSession(msg.Sender)
	if senderSession == nil {
		err = errors.NotFound("sender not found")
		return
	}

	msg.ID = uuid.NewV4().String()
	msg.CreatedAt = time.Now().UnixMilli()

	buf, err := codec.Pack(packet, msg)
	if err != nil {
		return
	}
	receiverSession.Send(buf)
	senderSession.Send(buf)

	app.Save(receiverSession.ID, msg)
	app.Save(senderSession.ID, msg)

	return
}

func (app *MessageApplication) Save(sessionId string, msg *Message) {
	app.rmw.Lock()
	defer app.rmw.Unlock()
	items := app.messages[sessionId]
	if len(items) == 0 {
		items = make([]*Message, 0, 100)
	}
	items = append(items, msg)
	if total := len(items); total > 100 {
		items = items[total-100 : total]
	}
	app.messages[sessionId] = items

}

func NewMessageApplication(bucket *bucket.Bucket) *MessageApplication {
	return &MessageApplication{
		bucket:   bucket,
		messages: map[string][]*Message{},
	}
}
