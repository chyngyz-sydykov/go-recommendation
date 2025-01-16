package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/chyngyz-sydykov/go-recommendation/application"
)

func main() {

	fmt.Println("go-recommend/rabbitmq consumer started")
	app := application.InitializeApplication()
	defer app.ShutDown()
	app.Start()
	waitForShutdown()
}

func waitForShutdown() {
	// Create a channel to listen for OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-signalChan
	fmt.Printf("Received signal: %v. Exiting...\n", sig)
}
