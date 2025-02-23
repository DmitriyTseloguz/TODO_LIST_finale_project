package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	server "todo-list/lib/server"

	"github.com/jmoiron/sqlx"
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

	if !isAlreadyExist(DB_FILE) {
		os.Create(DB_FILE)
		var err = initializeDB(DB_FILE)

		if err != nil {
			log.Fatal(err)
		}
	}

	var mux = http.DefaultServeMux

	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.Handle("/api/nextdate", server.ApiHandlers["/api/nextdate"])

	fmt.Println("Start listening: http://localhost:" + PORT)

	var serverError = http.ListenAndServe(":"+PORT, mux)

	if serverError != nil {
		log.Fatal(serverError)
	}
}

func isAlreadyExist(file string) bool {
	dbFile := getApplicationFilePath(file)

	var _, err = os.Stat(dbFile)

	return err == nil
}

func initializeDB(name string) error {
	var err error
	var db *sqlx.DB

	db, err = sqlx.Connect("sqlite3", name)

	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(`
CREATE TABLE scheduler(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date DATE,
	title TEXT,
	comment TEXT,
	repeat VARCHAR(128)
);
CREATE INDEX scheduler_date ON scheduler(date);
	`)

	return err
}

func getApplicationFilePath(file string) string {
	var appPath, err = os.Executable()

	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(filepath.Dir(appPath), file)
}
