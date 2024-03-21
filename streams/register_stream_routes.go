package streams

import (
	"htmxgo/entities"
	"htmxgo/middlewares"
	"io"

	"github.com/gin-gonic/gin"
)

var ArenaQuizStream *entities.StreamEventEntity

func registerStreamChannel(c *gin.Context) {
	v, ok := c.Get("clientChan")

	if !ok {
		return
	}

	clientChan, ok := v.(entities.ClientChan)

	if !ok {
		return
	}

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-clientChan; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}

func RegisterStreamRoutes(router *gin.Engine) {
	router.GET("/event-stream", middlewares.SSEHeaderMiddleware, ArenaQuizStream.ServeHTTP(), registerStreamChannel)
}
