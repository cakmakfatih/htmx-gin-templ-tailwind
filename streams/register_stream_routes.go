package streams

import (
	"fmt"
	"htmxgo/entities"
	"htmxgo/middlewares"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

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
	streamEventEntity := NewStreamServer()

	duration := 15 * time.Minute
	endTime := time.Now().Add(duration)

	go func() {
		for {
			time.Sleep(time.Second * 1)

			remaining := time.Until(endTime)
			minutes := remaining / time.Minute
			seconds := (remaining % time.Minute) / time.Second

			timeString := fmt.Sprintf("%02d:%02d", minutes, seconds)

			streamEventEntity.Message <- timeString
		}
	}()

	router.GET("/event-stream", middlewares.SSEHeaderMiddleware, streamEventEntity.ServeHTTP(), registerStreamChannel)
}
