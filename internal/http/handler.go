package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}
func (handler *Handler) Install(router *gin.Engine) {
	router.LoadHTMLFiles("web/index.html", "web/dark.html")
	router.StaticFS("/web/static", http.Dir("./web/static"))
	router.GET("/web/index", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.GET("/web/dark", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "dark.html", gin.H{})
	})
}
