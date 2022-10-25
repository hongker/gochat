package http

import "github.com/gin-gonic/gin"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}
func (handler *Handler) Install(router *gin.Engine) {

}
