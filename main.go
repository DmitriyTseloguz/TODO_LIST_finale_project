package main

import (
	"github.com/joho/godotenv"
	"os"

	"net/http"
)

const webDir = "./web"

func init(){
	godotenv.Load()
}

func main(){
	var PORT = os.Getenv("TODO_PORT")
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	http.ListenAndServe(":" + PORT, nil)
}