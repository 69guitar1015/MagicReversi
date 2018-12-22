package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func reserveFinalizeWhenExited(f interface {
	Finalize() error
}) {
	// Finalizing processing when termination signal comes
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		<-signalChan
		fmt.Println("Terminated...")
		f.Finalize()
		os.Exit(1)
	}()
}
