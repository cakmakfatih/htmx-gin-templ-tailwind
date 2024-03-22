package tasks

import (
	"fmt"
	"htmxgo/models"
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
	arenaQuizDuration := 5 * time.Minute

	startTime := time.Now().Truncate(arenaQuizIntervalDuration)
	endTime := startTime.Add(arenaQuizDuration)

	streams.ArenaQuizStream = streams.NewStreamServer()

	_, err = scheduler.NewJob(
		gocron.DurationJob(
			arenaQuizIntervalDuration,
		),
		gocron.NewTask(
			func(st *time.Time, et *time.Time) {
				*st = time.Now().Truncate(time.Minute).Add(arenaQuizIntervalDuration)
				*et = st.Add(arenaQuizDuration)

				models.GameModel{
					StartTime:  st.UTC(),
					FinishTime: et.UTC(),
				}.Insert()
			},
			&startTime,
			&endTime,
		),
		gocron.WithSingletonMode(1),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)

	if err != nil {
		log.Fatal(err)
	}

	_, err = scheduler.NewJob(
		gocron.DurationJob(
			time.Second,
		),
		gocron.NewTask(
			func(t *time.Time) {
				remaining := time.Until(*t)

				minutes := remaining / time.Minute
				seconds := (remaining % time.Minute) / time.Second

				timeString := fmt.Sprintf("%02d:%02d", minutes, seconds)

				streams.ArenaQuizStream.Message <- timeString
			},
			&startTime,
		),
		gocron.WithSingletonMode(1),
	)

	if err != nil {
		log.Fatal(err)
	}

	scheduler.Start()
}
