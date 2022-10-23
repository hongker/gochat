package application

import (
	"context"
	"sync"
)

type SessionApplication struct {
	rmu      sync.RWMutex
	sessions map[string][]*Session
}

type Session struct {
	ID    string
	Title string
	Last  *Message
}

func (app *SessionApplication) GetSessionList(ctx context.Context, uid string) ([]*Session, error) {
	app.rmu.RLock()
	items := app.sessions[uid]
	app.rmu.RUnlock()
	return items, nil
}

func (app *SessionApplication) SaveSession(ctx context.Context, uid string, session *Session) (err error) {
	app.rmu.Lock()
	defer app.rmu.Unlock()
	items := app.sessions[uid]
	if len(items) == 0 {
		app.sessions[uid] = []*Session{session}
		return
	}

	for _, item := range items {
		if item.ID == session.ID {
			item.Title = session.Title
			item.Last = session.Last
			break
		}
	}
	app.sessions[uid] = items

	return
}

func NewSessionApplication() *SessionApplication {
	return &SessionApplication{
		sessions: map[string][]*Session{},
	}
}
