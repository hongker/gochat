package application

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"gochat/pkg/cmap"
	"time"
)

type UserApplication struct {
	collection *cmap.Container[string, *User]
}

func NewUserApplication() *UserApplication {
	return &UserApplication{
		collection: cmap.NewContainer[string, *User](),
	}
}

type User struct {
	ID        string
	Name      string
	CreatedAt int64
}

// Auth represents user authentication
func (app *UserApplication) Auth(ctx context.Context, user *User) error {
	user.ID = uuid.NewV4().String()
	user.CreatedAt = time.Now().Unix()

	app.collection.Set(user.ID, user)
	return nil
}
