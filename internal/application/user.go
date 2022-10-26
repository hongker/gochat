package application

import (
	"context"
	"gochat/pkg/cmap"
	"gochat/pkg/gen"
	"time"
)

type UserApplication struct {
	collection *cmap.Container[string, *User]
	generator  gen.IDGenerator
}

func NewUserApplication() *UserApplication {
	return &UserApplication{
		collection: cmap.NewContainer[string, *User](),
		generator:  gen.NewSnowFlakeGenerator(),
	}
}

type User struct {
	ID        string
	Name      string
	CreatedAt int64
}

// Auth represents user authentication
func (app *UserApplication) Auth(ctx context.Context, user *User) error {
	user.ID = app.generator.Generate()
	user.CreatedAt = time.Now().Unix()

	app.collection.Set(user.ID, user)
	return nil
}
