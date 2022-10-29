package socket

import (
	"github.com/ebar-go/znet"
	"gochat/internal/bucket"
	"gochat/internal/domain/dto"
)

// connect represents user connect action
func (handler *Handler) connect(ctx *znet.Context, req *dto.ConnectRequest) (resp *dto.ConnectResponse, err error) {

	handler.setCurrentUser(ctx, req.UID)

	handler.bucket.AddSession(bucket.NewSession(req.UID, ctx.Conn()))

	return
}

// updateProfile updates the profile of user
func (handler *Handler) updateProfile(ctx *znet.Context, req *dto.UserUpdateRequest) (resp *dto.UserUpdateResponse, err error) {
	err = handler.userApp.Update(ctx, handler.currentUser(ctx), req)
	return
}

// findProfile returns the user profile
func (handler *Handler) findProfile(ctx *znet.Context, req *dto.IDRequest) (resp *dto.UserResponse, err error) {
	user, err := handler.userApp.Get(ctx, req.ID)
	if err != nil {
		return
	}

	resp = &dto.UserResponse{
		ID:        req.ID,
		Name:      user.Name,
		Avatar:    user.Avatar,
		Email:     user.Email,
		Location:  user.Location,
		CreatedAt: user.CreatedAt,
	}
	return
}

func (handler *Handler) getContacts(ctx *znet.Context, req *dto.ContactQueryRequest) (resp *dto.ContactQueryResponse, err error) {
	contacts, err := handler.userApp.GetContacts(ctx, handler.currentUser(ctx))
	if err != nil {
		return
	}
	resp = &dto.ContactQueryResponse{Items: make([]dto.User, 0)}
	for _, user := range contacts {
		resp.Items = append(resp.Items, dto.User{
			ID:        user.ID,
			Name:      user.Name,
			Avatar:    user.Avatar,
			Email:     user.Email,
			Location:  user.Location,
			CreatedAt: user.CreatedAt,
		})
	}
	return
}

const (
	userIdKey = "uid"
)

// setCurrentUser sets the current user id
func (handler *Handler) setCurrentUser(ctx *znet.Context, uid string) {
	ctx.Conn().Property().Set(userIdKey, uid)
}

// currentUser returns the current user id
func (handler *Handler) currentUser(ctx *znet.Context) string {
	uid, exist := ctx.Conn().Property().Get(userIdKey)
	if !exist {
		return ""
	}
	return uid.(string)
}
