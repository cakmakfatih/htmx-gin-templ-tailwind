package tasks

import (
	"fmt"
	"htmxgo/streams"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func InitScheduler() {
	scheduler, err := gocron.NewScheduler()

	if err != nil {
		log.Fatal(err)
	}

	arenaQuizIntervalDuration := 5 * time.Minute
	endTime := time.Now().Add(arenaQuizIntervalDuration)

	streams.ArenaQuizStream = streams.NewStreamServer()

	_, err = scheduler.NewJob(
		gocron.DurationJob(
			arenaQuizIntervalDuration,
		),
		gocron.NewTask(
			func() {
				endTime = time.Now().Add(arenaQuizIntervalDuration)
			},
		),
		gocron.WithSingletonMode(1),
	)

	if err != nil {
		log.Fatal(err)
	}

	_, err = scheduler.NewJob(
		gocron.DurationJob(
			time.Second,
		),
		gocron.NewTask(
			func() {
				remaining := time.Until(endTime)
				minutes := remaining / time.Minute
				seconds := (remaining % time.Minute) / time.Second

				timeString := fmt.Sprintf("%02d:%02d", minutes, seconds)

				streams.ArenaQuizStream.Message <- timeString
			},
		),
		gocron.WithSingletonMode(1),
	)

	if err != nil {
		log.Fatal(err)
	}

	scheduler.Start()
}
