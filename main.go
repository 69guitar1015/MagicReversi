package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	mrbt "github.com/69guitar1015/MagicReversi/mrbluetooth"
	"github.com/69guitar1015/MagicReversi/mrmiddle"
	"github.com/comail/colog"
)

const (
	TIMING_POLL = 50
)

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
			time.Sleep(req.Interval)
		}
	}

	getBoardHandle := func() mrbt.Payload {
		board, err := m.ReadBoard()

		if err != nil {
			log.Println("error: ", err)
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

func flip_debug() {
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

func read_debug() {
	m, err := mrmiddle.NewMrMiddle()
	checkError(err)

	defer m.Finalize()
	reserveFinalizeWhenExited(m)

	err = m.Init()
	checkError(err)

	for {
		b, err := m.ReadBoard()
		checkError(err)

		b.Print()
		time.Sleep(500 * time.Millisecond)
	}
}

func json_debug() {
	m, err := mrmiddle.NewMrMiddle()
	checkError(err)

	defer m.Finalize()
	reserveFinalizeWhenExited(m)

	err = m.Init()
	checkError(err)

	board, err := m.ReadBoard()

	var payload mrbt.Payload
	if err != nil {
		log.Println("error: ", err)
		payload = &ErrorPayload{"Error!"}
	} else {
		payload = &GetBoardPayload{board: board}
	}

	c := payload.Compose()
	fmt.Println(string(c))
}

func help() {
	fmt.Println("Valid command is {start, debug}")
}

func main() {
	colog.Register()
	colog.SetFlags(log.Ldate | log.Lshortfile)
	colog.SetMinLevel(colog.LDebug)

	if len(os.Args) != 2 {
		help()
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "start":
		start()
	case "debug":
		fmt.Println(rand.Int63())

		colog.SetMinLevel(colog.LDebug)
		fmt.Println("debug: Debug mode!")
		json_debug()
	default:
		help()
	}
}
