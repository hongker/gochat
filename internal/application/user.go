package application

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

type UserApplication struct {
	rmu   sync.RWMutex
	items map[string]*User
}

func NewUserApplication() *UserApplication {
	return &UserApplication{}
}

type User struct {
	ID        string
	Name      string
	CreatedAt int64
}

func (app *UserApplication) Login(ctx context.Context, user *User) error {
	user.ID = uuid.NewV4().String()
	user.CreatedAt = time.Now().Unix()

	app.rmu.Lock()
	app.items[user.ID] = user
	app.rmu.Unlock()
	return nil
}
