package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	planner "todo-list/lib"
	"todo-list/lib/model"

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

	mux.Handle("/api/nextdate", http.HandlerFunc(
		func(response http.ResponseWriter, request *http.Request) {
			// if request.Method != http.MethodGet {
			// 	response.WriteHeader(http.StatusMethodNotAllowed)
			// 	return
			// }

			var nowParameter = request.URL.Query().Get("now")
			var dateParameter = request.URL.Query().Get("date")
			var repeatParameter = request.URL.Query().Get("repeat")

			var baseOnDate, nowParseError = time.Parse("20060102", nowParameter)
			var _, dateParseError = time.Parse("20060102", dateParameter)

			if dateParseError != nil {
				response.Write([]byte("Wrong date"))
				return
			}

			if dateParameter == "" {
				response.Write([]byte(""))
				return
			}

			var task = model.NewTask(0, "", dateParameter, "", repeatParameter)

			var rescheduler, err = planner.DefineRescheduler(task)

			if err != nil {
				response.Write([]byte(err.Error()))

				return
			}

			if nowParseError == nil {
				rescheduler.SetBaseOnDate(baseOnDate)
			}

			rescheduler.Reschedule(task)

			var nextTime = task.GetTime()

			response.Write([]byte(nextTime.Format("20060102")))
		}),
	)

	fmt.Println("Start listening: http://localhost:" + PORT)
	http.ListenAndServe(":"+PORT, mux)
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
