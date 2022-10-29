package http

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}
func (handler *Handler) Install(router *gin.Engine) {
	//router.LoadHTMLFiles("web/login.html", "web/index.html", "web/dark.html")
	//router.StaticFS("/web/static", http.Dir("./web/static"))
	//
	//router.GET("/web/login", func(ctx *gin.Context) {
	//	ctx.HTML(http.StatusOK, "login.html", gin.H{})
	//})
	//router.GET("/web/index", func(ctx *gin.Context) {
	//	ctx.HTML(http.StatusOK, "index.html", gin.H{})
	//})
	//router.GET("/web/dark", func(ctx *gin.Context) {
	//	ctx.HTML(http.StatusOK, "dark.html", gin.H{})
	//})

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
