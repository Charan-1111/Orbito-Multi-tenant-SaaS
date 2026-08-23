package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"task-manager-auth/internal/server"
)

func main() {
	app, err := server.NewApplication()
	if err != nil {
		log.Fatal(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	err = app.StartApplication(shutdown)
	if err != nil {
		log.Fatal(err)
	}
}
