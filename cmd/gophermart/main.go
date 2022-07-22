package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gophermarket/internal/server"
)

func main() {

	server.Run()

	waitChan := make(chan struct{})

	go func() {

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
		<-sigChan

		fmt.Printf("Service closed at %s. Goodbye!\n", time.Now().Format(time.RFC3339))

		close(waitChan)
	}()

	<-waitChan
}
