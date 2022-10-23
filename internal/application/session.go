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
}

func (app *SessionApplication) GetSessionList(ctx context.Context, uid string) ([]*Session, error) {
	app.rmu.RLock()
	items := app.sessions[uid]
	app.rmu.RUnlock()
	return items, nil
}

func NewSessionApplication() *SessionApplication {
	return &SessionApplication{
		sessions: map[string][]*Session{},
	}
}
