package core

import "os"

type AppConfig struct {
	Port      string
	Providers []string
}

var Config AppConfig

func InitConfig() {
	Config = AppConfig{
		Port:      os.Getenv("PORT"),
		Providers: []string{"discord", "google"},
	}
}
