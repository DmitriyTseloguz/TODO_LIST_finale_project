package main

import (
	"log"
	"net/http"
	"os"

	"todo-list/lib/database"
	"todo-list/lib/server"

	_ "github.com/mattn/go-sqlite3"

	"github.com/joho/godotenv"
)

const webDir = "./web"

func init() {
	godotenv.Load()
}

func main() {
	var PORT = os.Getenv("TODO_PORT")
	var DB_FILE = os.Getenv("TODO_DBFILE")

	db, err := database.NewDB(DB_FILE)

	if err != nil {
		log.Fatalf("Failed to load database: %v", err)
	}

	defer db.Close()

	if err := db.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	todoServer := server.New(webDir)

	log.Println("Start listening: http://localhost:" + PORT)

	if err := http.ListenAndServe(":"+PORT, todoServer); err != nil {
		log.Fatal(err)
	}
}
