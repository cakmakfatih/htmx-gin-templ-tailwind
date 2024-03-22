package models

import (
	"htmxgo/core"
	"log"
	"time"
)

type GameModel struct {
	ID         string    `json:"-"`
	CreatedAt  time.Time `json:"-"`
	StartTime  time.Time `json:"start_time"`
	FinishTime time.Time `json:"finish_time"`
}

func (GameModel) Get(id string) (*GameModel, error) {
	var result GameModel
	err := core.DbClient.DB.From("games").Select("id").Single().Eq("id", id).Execute(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (g GameModel) Insert() (*GameModel, error) {
	var result []GameModel

	err := core.DbClient.DB.From("games").Insert(g).Execute(&result)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &result[0], nil
}
