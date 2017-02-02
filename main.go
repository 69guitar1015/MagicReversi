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

func checkError(err error, m *mrmiddle.MrMiddle) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	m, err := mrmiddle.NewMrMiddle()

	checkError(err, m)

	defer m.Finalize()

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

	checkError(err, m)

	g := mrsoft.NewGame(m)

	err = g.Start()

	checkError(err, m)
}
