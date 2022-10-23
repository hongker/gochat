package application

import (
	"context"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/znet/codec"
	uuid "github.com/satori/go.uuid"
	"gochat/internal/bucket"
	"time"
)

type MessageApplication struct {
	bucket *bucket.Bucket
}

type Message struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	Target      string `json:"target"`
	Sender      string `json:"sender"`
	CreatedAt   int64  `json:"created_at"`
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

	return
}

func NewMessageApplication(bucket *bucket.Bucket) *MessageApplication {
	return &MessageApplication{
		bucket: bucket,
	}
}
