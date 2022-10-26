package application

import (
	"context"
	"github.com/ebar-go/ego/errors"
	"gochat/internal/domain/dto"
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
	Avatar    string
	Sex       string
	Age       int
	CreatedAt int64
}

// Auth represents user authentication
func (app *UserApplication) Auth(ctx context.Context, user *User) error {
	user.ID = app.generator.Generate()
	user.CreatedAt = time.Now().Unix()

	app.collection.Set(user.ID, user)
	return nil
}

func (app *UserApplication) Get(ctx context.Context, uid string) (*User, error) {
	user, exist := app.collection.Get(uid)
	if !exist {
		return nil, errors.NotFound("user not found")
	}
	return user, nil
}

func (app *UserApplication) Update(ctx context.Context, uid string, req *dto.UserUpdateRequest) error {
	user, err := app.Get(ctx, uid)
	if err != nil {
		return err
	}

	user.Name = req.Name
	user.Avatar = req.Avatar
	user.Sex = req.Sex
	user.Age = req.Age

	app.collection.Set(user.ID, user)
	return nil
}
