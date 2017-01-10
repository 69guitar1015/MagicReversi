package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/69guitar1015/MagicReversi/mrmiddle"
)

// this reversi system manages the board status
// by using only 1, 0, -1 integer as black, none, white
const (
	BLACK = 1
	NONE  = 0
	WHITE = -1
	WALL  = 2
)

type state int

func (s *state) flip() {
	*s = -1 * *s
}

type player int

func (p player) enemy() player {
	return -1 * p
}

func (p player) color() state {
	return state(p)
}

type point [2]int

func (b point) equal(a point) bool {
	return a[0] == b[0] && a[1] == b[1]
}

type direction [2]int

type board [10][10]state

func (b *board) put(p point, c state) error {
	if b[p[1]][p[0]] != NONE {
		return errors.New("There is already put a stone")
	}

	b[p[1]][p[0]] = c
	return nil
}

func (b *board) flip(p point) {
	b[p[1]][p[0]].flip()
}

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
	record map[point]player
	// available points
	available map[point][]direction
}

func newGame(m middleware) (g *game) {
	g = &game{
		b: board{
			[10]state{WALL, WALL, WALL, WALL, WALL, WALL, WALL, WALL, WALL, WALL},
			[10]state{WALL, NONE, NONE, NONE, NONE, NONE, NONE, NONE, NONE, WALL},
			[10]state{WALL, NONE, NONE, NONE, NONE, NONE, NONE, NONE, NONE, WALL},
			[10]state{WALL, NONE, NONE, NONE, NONE, NONE, NONE, NONE, NONE, WALL},
			[10]state{WALL, NONE, NONE, NONE, WHITE, BLACK, NONE, NONE, NONE, WALL},
			[10]state{WALL, NONE, NONE, NONE, BLACK, WHITE, NONE, NONE, NONE, WALL},
			[10]state{WALL, NONE, NONE, NONE, NONE, NONE, NONE, NONE, NONE, WALL},
			[10]state{WALL, NONE, NONE, NONE, NONE, NONE, NONE, NONE, NONE, WALL},
			[10]state{WALL, NONE, NONE, NONE, NONE, NONE, NONE, NONE, NONE, WALL},
			[10]state{WALL, WALL, WALL, WALL, WALL, WALL, WALL, WALL, WALL, WALL},
		},
		crr:       BLACK,
		m:         m,
		record:    map[point]player{},
		available: map[point][]direction{},
	}

	return
}

func (g *game) start() (err error) {
	for {
		g.printBoard()

		g.setAvailable()

		if g.isFinish() {
			fmt.Println("Finish!")
			g.printSummary()
			return
		}

		// skip if game is not finished and there is not available points
		if len(g.available) == 0 {
			fmt.Println("skipping")
			g.crr = g.crr.enemy()
			continue
		}

		x, y := g.m.GetInput()
		p := point{x, y}

		err = g.put(p)
		if err != nil {
			return fmt.Errorf("Failed to put the stone: %s", err)
		}

		g.crr = g.crr.enemy()

	}
}

// seek available point
func (g *game) seekAvailable() map[point][]direction {
	available := map[point][]direction{}

	for y := 1; y <= 8; y++ {
		for x := 1; x <= 8; x++ {
			if g.b[y][x] != NONE {
				continue
			}

			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}

					if g.b[y+dy][x+dx] == g.crr.enemy().color() {
						dist := 2
						for {
							if g.b[y+dist*dy][x+dist*dx] == g.crr.color() {
								p := point{x, y}
								available[p] = append(available[p], direction{dx, dy})
								break
							} else if g.b[y+dist*dy][x+dist*dx] == g.crr.enemy().color() {
								dist++
							} else {
								break
							}
						}
					}
				}
			}
		}
	}

	return available
}

func (g *game) setAvailable() {
	g.available = g.seekAvailable()
}

// put a stone to (x, y) address on the board
func (g *game) put(p point) (err error) {
	// return error if the point is not available
	fmt.Println(g.available, p)
	if len(g.available[p]) == 0 {
		return errors.New("Can't put stones there")
	}

	err = g.b.put(p, g.crr.color())

	if err != nil {
		return
	}

	for _, d := range g.available[p] {
		dist := 1
		for {
			dp := point{p[0] + dist*d[0], p[1] + dist*d[1]}
			if g.b[dp[1]][dp[0]] == g.crr.color() {
				break
			} else if g.b[dp[1]][dp[0]] == g.crr.enemy().color() {
				g.b.flip(dp)
			} else {
				return errors.New("Can't available this direction")
			}

			dist++
		}
	}

	g.record[p] = g.crr

	return
}

// judge whether the game is finished
func (g *game) isFinish() bool {
	// if each player has no available points, game is over
	if len(g.available) == 0 {
		g.crr = g.crr.enemy()
		ava := g.seekAvailable()
		g.crr = g.crr.enemy()

		if len(ava) == 0 {
			return true
		}
	}

	return false
}

func (g *game) printBoard() {
	for i, row := range g.b {
		for j, v := range row {
			switch v {
			case BLACK:
				fmt.Printf("○\t")
			case WHITE:
				fmt.Printf("●\t")
			case NONE:
				fmt.Printf(" \t")
			case WALL:
				switch {
				case j == 0 || j == 9:
					fmt.Printf("%d\t", i)
				default:
					fmt.Printf("%d\t", j)
				}
			}
		}
		fmt.Printf("\n")
	}
}

// print the game summary
func (g *game) printSummary() {
	fmt.Println("this is SUMMARY")
}

func main() {
	m := mrmiddle.NewMrMiddle()

	g := newGame(m)

	err := g.start()

	if err != nil {
		log.Fatal(err)
	}
}
