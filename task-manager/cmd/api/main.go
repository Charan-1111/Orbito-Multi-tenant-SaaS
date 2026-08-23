package main

import (
	"log"
	"task-manager/internal/server"
)

func main() {
	app, err := server.NewApplication()
	if err != nil {
		log.Fatal(err)
	}

	err = app.StartApplication()
	if err != nil {
		log.Fatal(err)
	}
}
