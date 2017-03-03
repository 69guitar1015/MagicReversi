package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/69guitar1015/MagicReversi/mrmiddle"
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

	fmt.Println(err)

	// checkError(err, m)

	// for i := 0; i < 20; i++ {
	// 	err = m.GotThem(byte(12 * i))
	// 	time.Sleep(300 * time.Millisecond)
	// }
	// m.Finalize()

	// err = m.GotThem(200)
	// m.Finalize()
	// checkError(err, m)

	/////
	// i := 0
	// s := [][]int{[]int{3, 4}, []int{4, 4}, []int{5, 4}, []int{6, 4}, []int{6, 5}, []int{5, 5}, []int{4, 5}, []int{3, 5}}
	pd := mrmiddle.N

	var mode string
	fmt.Println("modeを入れてくれ[input / output]")
	fmt.Scan(&mode)

	switch mode {
	case "input":
		for {
			m.GetInput()
		}

	case "output":
		for {
			fmt.Println("x y 磁力の方向(0/1)\nを入力してくれ")
			var x, y, d int
			fmt.Scan(&x, &y, &d)

			if d == 0 {
				pd = mrmiddle.N
			} else {
				pd = mrmiddle.S
			}

			err := m.Flip(x, y, pd)
			log.Println(err)
		}
	}
}
