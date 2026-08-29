package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"task-manager-auth/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: could not load .env file:", err)
	}

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
