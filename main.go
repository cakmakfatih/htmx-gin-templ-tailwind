package main

import (
	"fmt"
	"htmxgo/api/v1"
	"htmxgo/core"
	"htmxgo/middlewares"
	"htmxgo/models"
	pages "htmxgo/views/pages"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func staticCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/assets/") {
			c.Header("Cache-Control", "public, max-age=86400")
		}

		c.Next()
	}
}

func setupRouter(router *gin.Engine) {
	err := router.SetTrustedProxies(nil)

	if err != nil {
		log.Fatal(err)
	}

	router.Use(staticCacheMiddleware())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Static("/assets", "./public")

	ginHtmlRenderer := router.HTMLRender
	router.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}
}

func registerRoutes(router *gin.Engine) {
	api.RegisterApi(router)

	router.GET("/", middlewares.UserMiddleware, func(c *gin.Context) {
		userModel := models.UserModel.UserFromContext(models.UserModel{}, c)

		r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, pages.Home(userModel))
		c.Render(http.StatusOK, r)
	})
}

func main() {
	core.LoadEnv()
	core.InitDb()

	router := gin.Default()

	setupRouter(router)
	registerRoutes(router)

	router.Run(fmt.Sprintf("0.0.0.0:%v", os.Getenv("PORT")))
}
