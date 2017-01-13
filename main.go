package main

import (
	"fmt"
	"log"
	"time"

	"gobot.io/x/gobot/platforms/intel-iot/edison"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	e := edison.NewAdaptor()

	err := e.Connect()

	checkError(err)

	f := 0
	i := 0
	out1 := 0
	out2 := 0
	for {
		v1, _ := e.DigitalRead("4")

		if v1 == 1 && f == 0 {
			out1 = i
			out2 = 1 - i
			f++
		} else if v1 == 1 && (f < 3) {
			out1 = 0
			out2 = 0
			f++
		} else if v1 == 0 {
			out1 = 0
			out2 = 0
			f = 0
		} else {
			i = 1 - i
			f = 0
		}

		e.DigitalWrite("12", byte(out1))
		e.DigitalWrite("13", byte(out2))

		fmt.Printf("val1: %d,\tout1: %d,\tout2 %d\n", v1, out1, out2)

		time.Sleep(500 * time.Millisecond)
	}
}
