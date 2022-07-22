package main

import (
	"context"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app/di"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	app, err := di.InitializeApp(ctx, cancel)
	check(err)

	workingDirectory, err := os.Getwd()
	check(err)

	app.Start(check, workingDirectory)
	<-interrupt()
	app.Shutdown(check)
}

func interrupt() chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	return interrupt
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
