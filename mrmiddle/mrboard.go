package mrmiddle

import (
	"log"
	"time"

	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

// MrMiddle is Magic Reversi's middle ware object
type MrMiddle struct {
	e *edison.EdisonAdaptor
}

// NewMrMiddle returns MrMiddle instance
func NewMrMiddle() (mm *MrMiddle) { // Edison Adapter
	mm = &MrMiddle{}
	mm.e = edison.NewEdisonAdaptor("edison")
	mm.e.Connect()

	mm.Init()

	return
}

type row [8]bool

func (r row) reversed() (reversed row) {
	for i, v := range r {
		reversed[len(reversed)-i-1] = v
	}

	return
}

func checkError(err error) {
	if err != nil {
		// log.Fatal(err)
		log.Print("Error : ")
		log.Println(err)
	}
}

// convert from byte object to boolean array
func byte2Row(b byte) (r row) {
	for i := 0; i < 8; i++ {
		r[i] = (b&0x01 == 0x01)
		b >>= 1
	}

	return
}

// take y and returns the Expander's address and gpio from [GPIOA, GPIOB]
func y2AddrAndGpio(y int) (addr int, gpio int) {
	// Expander address
	addr = EXOA[int(y/2)]

	// Use GPIO B when y is odd number
	gpio = GPIOA
	if y%2 != 0 {
		gpio = GPIOB
	}

	return
}

// Init is initialization function of MrMiddle
func (mm *MrMiddle) Init() {
	log.Println("initialize circuit...")
	var err error

	for _, addr := range EXIA {
		mm.e.I2cStart(addr)

		//　Initialize IOCON
		err = mm.e.I2cWrite(addr, []byte{IOCON, 0x00})
		checkError(err)

		// Initialize IODIR as read
		err = mm.e.I2cWrite(addr, []byte{IODIRA, 0x00})
		checkError(err)
		err = mm.e.I2cWrite(addr, []byte{IODIRB, 0x00})
		checkError(err)
	}

	for _, addr := range EXOA {
		mm.e.I2cStart(addr)

		//　Initialize IOCON
		err = mm.e.I2cWrite(addr, []byte{IOCON, 0x00})
		checkError(err)

		// Initialize IODIR as read
		err = mm.e.I2cWrite(addr, []byte{IODIRA, 0xFF})
		checkError(err)
		err = mm.e.I2cWrite(addr, []byte{IODIRB, 0xFF})
		checkError(err)
	}

	mm.writeAllLow()
}

// read given y line
func (mm *MrMiddle) readLine(y int) (r row) {
	addr, gpio := y2AddrAndGpio(y)

	mm.e.I2cStart(addr)

	err := mm.e.I2cWrite(addr, []byte{byte(gpio)})
	checkError(err)

	data, err := mm.e.I2cRead(addr, 1)
	checkError(err)

	r = byte2Row(data[0])
	if y%2 == 0 {
		r = r.reversed()
	}

	return
}

// read both A and B bits of given address of the Expander
func (mm MrMiddle) readAB(addr int) (byteSet [2]row) {
	mm.e.I2cStart(addr)

	err := mm.e.I2cWrite(addr, []byte{GPIOA})
	checkError(err)

	data, err := mm.e.I2cRead(addr, 2)
	checkError(err)

	byteSet[0] = byte2Row(data[0]).reversed()
	byteSet[1] = byte2Row(data[1])

	return
}

// ReadWholeBoard returns whole board status
func (mm *MrMiddle) readWholeBoard() (byteSet [8]row) {
	for i, addr := range EXIA {
		data := mm.readAB(addr)

		byteSet[2*i] = data[0].reversed()
		byteSet[2*i+1] = data[1]
	}

	return
}

// GetInput waits until board state changes and return x, y
func (mm *MrMiddle) GetInput() (int, int) {
	old := mm.readWholeBoard()
	for {
		crr := mm.readWholeBoard()

		if crr != old {
			for i, r := range crr {
				for j, v := range r {
					// if current status is True and old status is False then reutrn (x, y)
					if v && !old[i][j] {
						return j + 1, i + 1
					}
				}
			}
		}

		time.Sleep(POLLTIME)
	}
}

// write byte data to designated line
func (mm *MrMiddle) writeByte(y int, v byte) {
	addr, gpio := y2AddrAndGpio(y - 1)

	mm.e.I2cStart(addr)

	data := []byte{byte(gpio), v}
	mm.e.I2cWrite(addr, data)
}

func (mm *MrMiddle) writeAllLow() {
	for y := 0; y < 8; y++ {
		mm.writeByte(y, 0x00)
	}
}

// HighWhile make (x, y) to High while ms[msec]
func (mm *MrMiddle) highWhile(x int, y int, ms time.Duration) {
	bits := byte(0x01 << uint(x-1))

	mm.writeByte(y, bits)

	time.Sleep(ms * time.Millisecond)

	mm.writeByte(y, 0x00)
}

// Flip flips a stone at (x, y)
func (mm *MrMiddle) Flip(x int, y int) {
	mm.highWhile(x, y, FLIPTIME)
}
