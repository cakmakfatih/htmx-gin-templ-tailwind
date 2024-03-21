package core

import "os"

type AppConfig struct {
	SupabaseUrl string
	SupabaseKey string
	Port        string
	Providers   []string
}

var Config AppConfig

func InitConfig() {
	Config = AppConfig{
		SupabaseUrl: os.Getenv("SUPABASE_URL"),
		SupabaseKey: os.Getenv("SUPABASE_KEY"),
		Port:        os.Getenv("PORT"),
		Providers:   []string{"discord", "google"},
	}
}
