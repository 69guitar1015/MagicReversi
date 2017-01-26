package main

import (
	"log"

	"github.com/69guitar1015/MagicReversi/mrmiddle"
	"github.com/69guitar1015/MagicReversi/mrsoft"
)

func checkError(err error, m *mrmiddle.MrMiddle) {
	if err != nil {
		m.Off()
		log.Fatal(err)
	}
}

func main() {
	m, err := mrmiddle.NewMrMiddle()

	checkError(err, m)

	err = m.Init()

	checkError(err, m)

	g := mrsoft.NewGame(m)

	err = g.Start()

	checkError(err, m)
}
