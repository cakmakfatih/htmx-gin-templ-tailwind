package core

import (
	"log"
	"os"

	supa "github.com/nedpals/supabase-go"
)

var DbClient *supa.Client

func InitDb() {
	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	DbClient = supa.CreateClient(supabaseUrl, supabaseKey)

	log.Println("Connected to Supabase")
}
