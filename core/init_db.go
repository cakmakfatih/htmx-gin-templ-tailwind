package core

import (
	"log"

	supa "github.com/nedpals/supabase-go"
)

var DbClient *supa.Client

func InitDb() {
	DbClient = supa.CreateClient(Config.SupabaseUrl, Config.SupabaseKey)

	log.Println("Connected to Supabase")
}
