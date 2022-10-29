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

var userCollection = cmap.NewContainer[string, *User]()

func NewUserApplication() *UserApplication {
	return &UserApplication{
		collection: userCollection,
		generator:  gen.NewSnowFlakeGenerator(),
	}
}

type User struct {
	ID        string
	Name      string
	Avatar    string
	Email     string
	Location  string
	Status    string
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
		return nil, errors.NotFound("user not found: %s", uid)
	}
	return user, nil
}

func (app *UserApplication) Update(ctx context.Context, uid string, req *dto.UserUpdateRequest) error {
	user, err := app.Get(ctx, uid)
	if err != nil {
		return err
	}

	//user.Name = req.Name
	user.Avatar = req.Avatar
	user.Location = req.Location
	user.Email = req.Email

	app.collection.Set(user.ID, user)
	return nil
}
