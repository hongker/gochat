package socket

import (
	"github.com/ebar-go/znet"
	"gochat/internal/application"
	"gochat/internal/bucket"
	"gochat/internal/domain/dto"
)

// login represents user login action
func (handler *Handler) login(ctx *znet.Context, req *dto.LoginRequest) (resp *dto.LoginResponse, err error) {
	user := &application.User{Name: req.Name}

	err = handler.userApp.Auth(ctx, user)
	if err != nil {
		return
	}

	handler.setCurrentUser(ctx, user.ID)

	handler.bucket.AddSession(bucket.NewSession(user.ID, ctx.Conn()))

	resp = &dto.LoginResponse{UID: user.ID}
	return
}

// updateProfile updates the profile of user
func (handler *Handler) updateProfile(ctx *znet.Context, req *dto.UserUpdateRequest) (resp *dto.UserUpdateRequest, err error) {
	err = handler.userApp.Update(ctx, handler.currentUser(ctx), req)
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
