package entities

import (
	"log"

	"github.com/gin-gonic/gin"
)

type StreamEventEntity struct {
	Message       chan string
	NewClients    chan chan string
	ClosedClients chan chan string
	TotalClients  map[chan string]bool
}

type ClientChan chan string

func (stream *StreamEventEntity) Listen() {
	for {
		select {
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(stream.TotalClients))

		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(stream.TotalClients))

		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.TotalClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

func (stream *StreamEventEntity) ServeHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientChan := make(ClientChan)

		stream.NewClients <- clientChan

		defer func() {
			stream.ClosedClients <- clientChan
		}()

		c.Set("clientChan", clientChan)

		c.Next()
	}
}
