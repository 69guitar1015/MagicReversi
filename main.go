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

func (s state) pole() mrmiddle.Pole {
	return mrmiddle.Pole(s)
}

func (s state) String() string {
	switch s {
	case BLACK:
		return "BLACK"
	case WHITE:
		return "WHITE"
	case NONE:
		return "NONE"
	case WALL:
		return "WALL"
	default:
		return "UNKNOWN"
	}
}

type player int

func (p player) enemy() player {
	return -1 * p
}

func (p player) color() state {
	return state(p)
}

func (p player) String() string {
	switch p {
	case BLACK:
		return "BLACK"
	case WHITE:
		return "WHITE"
	default:
		return "UNKNOWN"
	}
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
	Init() error
	GetInput() (int, int, error)
	Flip(int, int, mrmiddle.Pole) error
}

// single record of put history
type putRecord struct {
	pt    point
	pl    player
	flips []point
}

// game is main reversi game object
type game struct {
	// board object
	b board
	// current player
	crr player
	// middleware object
	m middleware
	// history of put stone
	history []putRecord
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
		history:   []putRecord{},
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

		x, y, err := g.m.GetInput()

		if err != nil {
			return fmt.Errorf("Failed to get input: %s", err)
		}

		p := point{x, y}

		if p[0] == -1 && p[1] == -1 {
			// undo when (x, y) == (-1, -1)
			err = g.undo()

			if err != nil {
				return fmt.Errorf("Failed to undo: %s", err)
			}

			continue
		}

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
	if len(g.available[p]) == 0 {
		return errors.New("Can't put stones there")
	}

	pr := putRecord{
		pt:    p,
		pl:    g.crr,
		flips: []point{},
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

				err = g.flip(dp, g.crr.color())

				if err != nil {
					return
				}

				// append to flip record
				pr.flips = append(pr.flips, dp)
			} else {
				return errors.New("Can't available this direction")
			}

			dist++
		}
	}

	g.history = append(g.history, pr)

	return
}

func (g *game) flip(p point, s state) (err error) {
	// simulational flip
	g.b.flip(p)
	// physical flip
	err = g.m.Flip(p[0], p[1], s.pole())

	return
}

func (g *game) undo() (err error) {
	record := g.history[len(g.history)-1]
	g.history = g.history[:len(g.history)-1]

	g.b[record.pt[1]][record.pt[0]] = NONE

	for i := range record.flips {
		// re-flip backwards
		p := record.flips[len(record.flips)-i-1]
		err = g.flip(p, record.pl.enemy().color())

		if err != nil {
			return fmt.Errorf("Failed to flip: %s", err)
		}
	}

	g.crr = record.pl

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
	fmt.Println("# SUMMARY ###########################################################")

	counts := map[state]int{BLACK: 0, WHITE: 0, NONE: 0}

	for i := 0; i <= 8; i++ {
		for j := 0; j <= 8; j++ {
			counts[g.b[i][j]]++
		}
	}

	switch {
	case counts[BLACK] == counts[WHITE]:
		fmt.Println("DRAW")
	default:
		fmt.Printf("%s PLAYER WINS!\n", player(BLACK))
	}

	fmt.Printf("NUMBER OF BLACK STONE:\t%2d\n", counts[BLACK])
	fmt.Printf("NUMBER OF WHITE STONE:\t%2d\n", counts[WHITE])
	fmt.Printf("NUMBER OF BLANK SPACE:\t%2d\n", counts[NONE])

	fmt.Printf("\n# KIFU\n")
	for i, record := range g.history {
		fmt.Printf("[%2d]\t(%d, %d)\t%s\t", i+1, record.pt[0], record.pt[1], record.pl)
		if (i+1)%3 == 0 {
			fmt.Printf("\n")
		}
	}

	fmt.Println("#####################################################################")
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	m, err := mrmiddle.NewMrMiddle()

	checkError(err)

	g := newGame(m)

	err = g.start()

	checkError(err)
}
