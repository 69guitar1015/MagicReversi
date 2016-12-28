package main

import (
	"log"
	"time"

	"github.com/69guitar1015/MagicReversi/mrmiddle"
)

func main() {
	b := mrmiddle.NewMrBoard()
	b.Init()

	for {
		ret := b.ReadWholeBoard()

		for _, bits := range ret {
			log.Println(bits)
		}

		time.Sleep(time.Second)
	}
}
