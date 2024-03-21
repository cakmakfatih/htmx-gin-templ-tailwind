package controllers

import (
	tabs "htmxgo/views/tabs"
	"net/http"

	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
	"github.com/gin-gonic/gin"
)

func renderGamesTab(c *gin.Context) {
	r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, tabs.GamesTab())
	c.Render(http.StatusOK, r)
}

func RegisterPartials(router *gin.Engine) {
	group := router.Group("/partials")

	group.GET("/tabs/games", renderGamesTab)
}
