package application

import (
	"context"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/znet/codec"
	uuid "github.com/satori/go.uuid"
	"gochat/internal/bucket"
	"gochat/internal/domain/constant"
	"gochat/internal/domain/dto"
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
	SessionID   string `json:"session_id"`
	SessionType string `json:"session_type"`
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

func (app *MessageApplication) GetLast(ctx context.Context, sessionID string) (item *Message) {
	app.rmw.RLock()
	items := app.messages[sessionID]
	app.rmw.RUnlock()
	if len(items) > 0 {
		item = items[len(items)-1]
	}
	return
}

func (app *MessageApplication) Send(ctx context.Context, sessionId string, sender, receiver *User, req *dto.MessageSendRequest, packet codec.Codec) (msg *Message, err error) {
	receiverSession := app.bucket.GetSession(req.Target)
	if receiverSession == nil {
		err = errors.NotFound("receiver not found")
		return
	}
	senderSession := app.bucket.GetSession(sender.ID)
	if senderSession == nil {
		err = errors.NotFound("sender not found")
		return
	}

	msg = &Message{
		ID:          uuid.NewV4().String(),
		SessionID:   sessionId,
		SessionType: constant.SessionTypeUser,
		Content:     req.Content,
		ContentType: req.ContentType,
		Target:      req.Target,
		Sender:      sender.ID,
		CreatedAt:   time.Now().UnixMilli(),
	}

	//app.Save(receiverSession.ID, msg)
	app.Save(msg.SessionID, msg)

	receiverBuf, err := packet.Pack(dto.Message{
		ID:           msg.ID,
		SessionID:    msg.SessionID,
		SessionTitle: sender.Name,
		Content:      msg.Content,
		ContentType:  msg.ContentType,
		CreatedAt:    msg.CreatedAt,
		Sender:       dto.User{ID: sender.ID, Name: sender.Name, Avatar: sender.Avatar},
	})
	if err == nil {
		receiverSession.Send(receiverBuf)
	}

	senderBuf, err := packet.Pack(dto.Message{
		ID:           msg.ID,
		SessionID:    msg.SessionID,
		SessionTitle: receiver.Name,
		Content:      msg.Content,
		ContentType:  msg.ContentType,
		CreatedAt:    msg.CreatedAt,
		Sender:       dto.User{ID: sender.ID, Name: sender.Name, Avatar: sender.Avatar},
	})
	if err == nil {
		senderSession.Send(senderBuf)
	}

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
