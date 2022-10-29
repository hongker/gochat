package http

import (
	"context"
	"github.com/ebar-go/ego/errors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gochat/internal/application"
	"gochat/internal/domain/dto"
	"os"
	"strings"
)

type Handler struct {
	userApp *application.UserApplication
}

func NewHandler() *Handler {
	return &Handler{
		userApp: application.NewUserApplication(),
	}
}
func (handler *Handler) Install(router *gin.Engine) {
	router.POST("/user/auth", convertAction[dto.LoginRequest, dto.LoginResponse](handler.login))
	router.Use(static.Serve("/", static.LocalFile("app/dist", true)))
	router.NoRoute(func(ctx *gin.Context) {
		accept := ctx.Request.Header.Get("Accept")
		flag := strings.Contains(accept, "text/html")
		if flag {
			content, err := os.ReadFile("app/dist/index.html")
			if err != nil {
				ctx.AbortWithStatus(404)
			}
			ctx.Writer.WriteHeader(200)
			ctx.Writer.Header().Add("Accept", "text/html")
			ctx.Writer.Write(content)
			ctx.Writer.Flush()
		}
	})

}

func (handler *Handler) login(ctx context.Context, req *dto.LoginRequest) (resp *dto.LoginResponse, err error) {
	user := &application.User{Name: req.Name}

	err = handler.userApp.Auth(ctx, user)
	if err != nil {
		return
	}
	resp = &dto.LoginResponse{UID: user.ID, Token: uuid.NewV4().String()}
	return
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func convertAction[Request, Response any](action func(ctx context.Context, req *Request) (*Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := new(Request)
		if err := ctx.ShouldBind(req); err != nil {
			ctx.JSON(200, Result{Code: 1001, Msg: "invalid param"})
			ctx.Abort()
			return
		}
		resp, err := action(ctx, req)

		result := Result{Data: resp}
		if err != nil {
			e := errors.Convert(err)
			result.Code = e.Code()
			result.Msg = e.Message()
		}

		ctx.JSON(200, result)

	}
}
