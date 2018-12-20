package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	mrbt "github.com/69guitar1015/MagicReversi/mrbluetooth"
	"github.com/69guitar1015/MagicReversi/mrmiddle"
)

const (
	TIMING_POLL = 50
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

type GetBoardPayload struct {
	board mrmiddle.Board `json:"board"`
}

func (p *GetBoardPayload) Compose() []byte {
	jsonBytes, err := json.Marshal(*p)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return nil
	}
	return jsonBytes
}

type NotifyPayload struct {
	x uint8 `json:"x"`
	y uint8 `json:"y"`
}

func (p *NotifyPayload) Compose() []byte {
	jsonBytes, err := json.Marshal(*p)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return nil
	}
	return jsonBytes
}

type ErrorPayload struct {
	err string `json:error`
}

func (p *ErrorPayload) Compose() []byte {
	jsonBytes, err := json.Marshal(*p)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return nil
	}
	return jsonBytes
}

func start() {
	m, err := mrmiddle.NewMrMiddle()
	checkError(err)

	defer m.Finalize()
	reserveFinalizeWhenExited(m)

	err = m.Init()
	checkError(err)

	flipHandle := func(req mrbt.FlipRequest) {
		for _, r := range req.Seq {
			m.Flip(r.Y, r.X, mrmiddle.Pole(r.Pole))
			time.Sleep(time.Duration(req.Interval) * time.Millisecond)
		}
	}

	getBoardHandle := func() mrbt.Payload {
		board, err := m.ReadBoard()

		if err != nil {
			log.Println(err)
			return &ErrorPayload{"Error!"}
		}

		payload := GetBoardPayload{board: board}
		return &payload
	}

	notify_chan := make(chan mrbt.Payload, 64)
	defer close(notify_chan)

	// Do notify
	go func() {
		// TODO: チャタリング制御
		old, err := m.ReadBoard()
		if err != nil {
			log.Println(err)
			return
		}

		for {
			// wait
			time.Sleep(TIMING_POLL * time.Millisecond)

			crr, err := m.ReadBoard()
			if err != nil {
				continue
			}

			for i := range crr {
				for j := range crr[i] {
					if old[i][j] == 0 && crr[i][j] == 1 {
						notify_chan <- &NotifyPayload{x: uint8(j), y: uint8(i)}
					}
				}
			}
		}
	}()

	bt := mrbt.NewMrBluetooth()
	bt.Launch(flipHandle, getBoardHandle, notify_chan)

	select {}
}

func debug() {
	fmt.Println("debug mode!")
	m, err := mrmiddle.NewMrMiddle()
	checkError(err)

	defer m.Finalize()
	reserveFinalizeWhenExited(m)

	err = m.Init()
	checkError(err)

	// read input
	stdin := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("i : ")
		stdin.Scan()
		i, _ := strconv.Atoi(stdin.Text())

		fmt.Printf("j : ")
		stdin.Scan()
		j, _ := strconv.Atoi(stdin.Text())

		fmt.Printf("Pole : ")
		stdin.Scan()
		text := stdin.Text()

		switch text {
		case "N":
			m.Flip(uint8(i), uint8(j), mrmiddle.N)
		case "S":
			m.Flip(uint8(i), uint8(j), mrmiddle.S)
		default:
			fmt.Println("Invalid input: ", text)
		}
	}
}

func help() {
	fmt.Println("Valid command is {start, debug}")
}

func main() {
	if len(os.Args) != 2 {
		help()
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "start":
		start()
	case "debug":
		debug()
	default:
		help()
	}
}
