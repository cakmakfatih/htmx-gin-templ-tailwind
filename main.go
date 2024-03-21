package main

import (
	"fmt"
	"htmxgo/api/v1"
	"htmxgo/core"
	"htmxgo/middlewares"
	"htmxgo/models"
	pages "htmxgo/views/pages"
	tabs "htmxgo/views/tabs"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func staticCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasSuffix(c.Request.URL.Path, "/assets/") {
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

	router.GET("/partials/tabs/games", func(c *gin.Context) {
		time.Sleep(time.Second * 3)

		r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, tabs.GamesTab())
		c.Render(http.StatusOK, r)
	})
}

func main() {
	core.LoadEnv()
	core.InitConfig()
	core.InitDb()

	router := gin.Default()

	setupRouter(router)
	registerRoutes(router)

	router.Run(fmt.Sprintf("0.0.0.0:%v", core.Config.Port))
}
