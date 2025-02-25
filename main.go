package main

import (
	"fmt"
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

	var db, databaseError = database.NewDB(DB_FILE)

	if databaseError != nil {
		log.Fatalf("Failed to load database: %v", databaseError)
	}

	defer db.Close()

	if err := db.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	var mux = http.DefaultServeMux

	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	for rout, handler := range server.ApiHandlers {
		mux.Handle(rout, handler)
	}

	fmt.Println("Start listening: http://localhost:" + PORT)

	var serverError = http.ListenAndServe(":"+PORT, mux)

	if serverError != nil {
		log.Fatal(serverError)
	}
}
