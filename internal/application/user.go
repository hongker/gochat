package application

import (
	"context"
	"github.com/ebar-go/ego/errors"
	"gochat/internal/domain/dto"
	"gochat/internal/domain/types"
	"gochat/pkg/cmap"
	"gochat/pkg/uuid"
	"time"
)

type UserApplication struct {
	collection *cmap.Container[string, *types.User]
	generator  uuid.IDGenerator
	contacts   *cmap.Container[string, []string]
}

var userCollection = cmap.NewContainer[string, *types.User]()
var contacts = cmap.NewContainer[string, []string]()

func NewUserApplication() *UserApplication {
	return &UserApplication{
		collection: userCollection,
		generator:  uuid.NewSnowFlakeGenerator(),
		contacts:   contacts,
	}
}

// Auth represents user authentication
func (app *UserApplication) Auth(ctx context.Context, user *types.User) error {
	user.ID = app.generator.Generate()
	user.CreatedAt = time.Now().Unix()

	app.collection.Set(user.ID, user)
	return nil
}

func (app *UserApplication) FindByEmail(ctx context.Context, email string) (user *types.User) {
	app.collection.Iterator(func(key string, val *types.User) {
		if val.Email == email {
			user = val
		}
	})

	return
}

func (app *UserApplication) Get(ctx context.Context, uid string) (*types.User, error) {
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

func (app *UserApplication) GetContacts(ctx context.Context, uid string) (items []*types.User, err error) {
	items = make([]*types.User, 0, 64)
	ids, exist := app.contacts.Get(uid)
	if !exist {
		return
	}

	for _, id := range ids {
		user, exist := app.collection.Get(id)
		if !exist {
			continue
		}
		items = append(items, user)
	}

	return
}

func inArray[T comparable](items []T, item T) bool {
	for _, t := range items {
		if t == item {
			return true
		}
	}

	return false
}
func (app *UserApplication) SaveContact(ctx context.Context, uid string, targetId string) {
	ids, exist := app.contacts.Get(uid)
	if !exist {
		ids = []string{targetId}
	} else {
		if inArray[string](ids, targetId) {
			return
		}
		ids = append(ids, targetId)
	}

	app.contacts.Set(uid, ids)
}
