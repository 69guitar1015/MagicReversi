package main

import (
	"log"

	"github.com/69guitar1015/MagicReversi/mrmiddle"
)

// this reversi system manages the board status
// by using only 1, 0, -1 integer as black, none, white
const (
	BLACK = 1
	NONE  = 0
	WHITE = -1
)

type state int
type player int
type point [2]int
type board [8][8]state

// reversi middleware interface
type middleware interface {
	GetInput() (int, int)
	Flip(int, int)
}

// game is main reversi game object
type game struct {
	// board object
	b board
	// current player
	crr player
	// middleware object
	m middleware
	// history of putting stone
	record []point
}

func newGame(m middleware) (g *game) {
	g = &game{
		b: board{
			[8]state{NONE, NONE, NONE, NONE, NONE, NONE, NONE, NONE},
			[8]state{NONE, NONE, NONE, NONE, NONE, NONE, NONE, NONE},
			[8]state{NONE, NONE, NONE, NONE, NONE, NONE, NONE, NONE},
			[8]state{NONE, NONE, NONE, BLACK, WHITE, NONE, NONE, NONE},
			[8]state{NONE, NONE, NONE, WHITE, BLACK, NONE, NONE, NONE},
			[8]state{NONE, NONE, NONE, NONE, NONE, NONE, NONE, NONE},
			[8]state{NONE, NONE, NONE, NONE, NONE, NONE, NONE, NONE},
			[8]state{NONE, NONE, NONE, NONE, NONE, NONE, NONE, NONE},
		},
		crr: BLACK,
		m:   m,
	}

	return
}

func (g *game) start() {
	for {
		x, y := g.m.GetInput()

		err := g.put(point{x, y})
		if err != nil {
			log.Fatal(err)
		}

		if g.isFinish() {
			log.Println("Finish!")
			g.printSummary()
			break
		}
	}
}

// put a stone to (x, y) address on the board
func (g *game) put(p point) (err error) {
	return
}

// judge whether the game is finished
func (g *game) isFinish() bool {
	return true
}

// print the game summary
func (g *game) printSummary() {

}

func main() {
	m := mrmiddle.NewMrMiddle()

	g := newGame(m)

	g.start()
}
