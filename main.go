package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/69guitar1015/MagicReversi/mrmiddle"
	"github.com/69guitar1015/MagicReversi/mrsoft"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func main() {
	m, err := mrmiddle.NewMrMiddle()
	checkError(err)

	defer m.Finalize()

	// Finalizing processing when termination signal comes
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-signalChan
		fmt.Println("Terminated...")
		m.Finalize()
		os.Exit(1)
	}()

	err = m.Init()
	checkError(err)

	g := mrsoft.NewGame(m)

	err = g.Start()
	checkError(err)
}
