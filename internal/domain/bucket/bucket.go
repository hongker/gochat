package bucket

import (
	"github.com/ebar-go/ego/utils/runtime"
	"gochat/pkg/cmap"
	"sync"
	"sync/atomic"
)

// Bucket represents a bucket for all connections
type Bucket struct {
	*Room
	channels *cmap.Container[string, *Channel]

	once sync.Once
	done chan struct{}

	workerNum  uint64
	queueSize  uint64
	queueCount uint64
	queues     []chan QueueItem
}

type QueueItem struct {
	Channel *Channel
	Msg     []byte
}

func (bucket *Bucket) AddChannel(id string) *Channel {
	channel := NewChannel(id)
	bucket.channels.Set(id, channel)
	return channel
}
func (bucket *Bucket) RemoveChannel(channel *Channel) {
	bucket.channels.Del(channel.ID)
}
func (bucket *Bucket) GetChannel(id string) *Channel {
	channel, _ := bucket.channels.Get(id)
	return channel
}

func (bucket *Bucket) SubscribeChannel(channel *Channel, sessions ...*Session) {
	if channel == nil {
		return
	}
	for _, session := range sessions {
		channel.AddSession(session)
		session.Channels = append(session.Channels, channel.ID)
	}

}

func (bucket *Bucket) UnsubscribeChannel(channel *Channel, sessions ...*Session) {
	if channel == nil {
		return
	}
	for _, session := range sessions {
		channel.RemoveSession(session)
	}
}

func (bucket *Bucket) BroadcastChannel(channel *Channel, msg []byte) {
	if channel == nil {
		return
	}
	num := atomic.AddUint64(&bucket.workerNum, 1) % bucket.queueCount
	select {
	case bucket.queues[num] <- QueueItem{Channel: channel, Msg: msg}:
	default:
	}

}

func (bucket *Bucket) Stop() {
	bucket.once.Do(func() {
		close(bucket.done)
	})
}

func (bucket *Bucket) start() {
	for i := 0; i < int(bucket.queueCount); i++ {
		bucket.queues[i] = make(chan QueueItem, bucket.queueSize)
		go func(idx int) {
			defer runtime.HandleCrash()
			bucket.polling(bucket.done, bucket.queues[idx])
		}(i)
	}
}

func (bucket *Bucket) polling(done <-chan struct{}, queue chan QueueItem) {
	for {
		select {
		case <-done:
			return
		case item, ok := <-queue:
			if !ok {
				return
			}

			item.Channel.Broadcast(item.Msg)
		}
	}

}

func NewBucket(options Options) *Bucket {
	bucket := &Bucket{
		channels:   cmap.NewContainer[string, *Channel](),
		Room:       NewRoom(),
		queues:     make([]chan QueueItem, options.QueueCount),
		queueSize:  options.QueueSize,
		queueCount: options.QueueCount,
	}
	bucket.start()
	return bucket
}

type Options struct {
	QueueCount uint64
	QueueSize  uint64
}
