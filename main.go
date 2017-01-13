package main

import (
	"log"

	"github.com/69guitar1015/MagicReversi/mrmiddle"
	"github.com/69guitar1015/MagicReversi/mrsoft"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	m, err := mrmiddle.NewMrMiddle()

	checkError(err)

	err = m.Init()

	checkError(err)

	g := mrsoft.NewGame(m)

	err = g.Start()

	checkError(err)
}
