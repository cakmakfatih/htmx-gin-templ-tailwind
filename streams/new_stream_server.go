package streams

import "htmxgo/entities"

func NewStreamServer() (event *entities.StreamEventEntity) {
	event = &entities.StreamEventEntity{
		Message:       make(chan string),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
	}

	go event.Listen()

	return
}
