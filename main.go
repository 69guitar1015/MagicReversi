package main

import (
	"log"
	"time"

	"github.com/69guitar1015/MagicReversi/mrmiddle"
)

func main() {
	mm := mrmiddle.NewMrMiddle()
	mm.Init()

	for {
		ret := mm.ReadWholeBoard()

		for _, bits := range ret {
			log.Println(bits)
		}

		time.Sleep(time.Second)
	}
}
