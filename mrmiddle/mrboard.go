package mrmiddle

import (
	"log"
	"time"

	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

// MrBoard is Magic Reversi's board object
type MrBoard struct {
	e *edison.EdisonAdaptor
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

// NewMrBoard returns MrBoard instance
func NewMrBoard() (b *MrBoard) { // Edison Adapter
	b = &MrBoard{}
	b.e = edison.NewEdisonAdaptor("edison")
	b.e.Connect()

	return
}

// Init is initialization function of MrBoard
func (b *MrBoard) Init() {
	log.Println("initialize circuit...")
	var err error

	for _, addr := range EXIA {
		b.e.I2cStart(addr)

		//　Initialize IOCON
		err = b.e.I2cWrite(addr, []byte{IOCON, 0x00})
		checkError(err)

		// Initialize IODIR as read
		err = b.e.I2cWrite(addr, []byte{IODIRA, 0x00})
		checkError(err)
		err = b.e.I2cWrite(addr, []byte{IODIRB, 0x00})
		checkError(err)
	}

	for _, addr := range EXOA {
		b.e.I2cStart(addr)

		//　Initialize IOCON
		err = b.e.I2cWrite(addr, []byte{IOCON, 0x00})
		checkError(err)

		// Initialize IODIR as read
		err = b.e.I2cWrite(addr, []byte{IODIRA, 0xFF})
		checkError(err)
		err = b.e.I2cWrite(addr, []byte{IODIRB, 0xFF})
		checkError(err)
	}

	b.writeAllLow()
}

// read given y line
func (b *MrBoard) readLine(y int) (r row) {
	addr, gpio := y2AddrAndGpio(y)

	b.e.I2cStart(addr)

	err := b.e.I2cWrite(addr, []byte{byte(gpio)})
	checkError(err)

	data, err := b.e.I2cRead(addr, 1)
	checkError(err)

	r = byte2Row(data[0])
	if y%2 == 0 {
		r = r.reversed()
	}

	return
}

// read both A and B bits of given address of the Expander
func (b MrBoard) readAB(addr int) (byteSet [2]row) {
	b.e.I2cStart(addr)

	err := b.e.I2cWrite(addr, []byte{GPIOA})
	checkError(err)

	data, err := b.e.I2cRead(addr, 2)
	checkError(err)

	byteSet[0] = byte2Row(data[0]).reversed()
	byteSet[1] = byte2Row(data[1])

	return
}

// ReadWholeBoard returns whole board status
func (b *MrBoard) readWholeBoard() (byteSet [8]row) {
	for i, addr := range EXIA {
		data := b.readAB(addr)

		byteSet[2*i] = data[0].reversed()
		byteSet[2*i+1] = data[1]
	}

	return
}

// GetInput waits until board state changes and return x, y
func (b *MrBoard) GetInput() (int, int) {
	old := b.readWholeBoard()
	for {
		crr := b.readWholeBoard()

		if crr != old {
			for i, r := range crr {
				for j, v := range r {
					// if current status is True and old status is False then reutrn (x, y)
					if v && !old[i][j] {
						return j, i
					}
				}
			}
		}

		time.Sleep(POLLTIME)
	}
}

// write byte data to designated line
func (b *MrBoard) writeByte(y int, v byte) {
	addr, gpio := y2AddrAndGpio(y)

	b.e.I2cStart(addr)

	data := []byte{byte(gpio), v}
	b.e.I2cWrite(addr, data)
}

func (b *MrBoard) writeAllLow() {
	for y := 0; y < 8; y++ {
		b.writeByte(y, 0x00)
	}
}

// HighWhile make (x, y) to High while ms[msec]
func (b *MrBoard) highWhile(x int, y int, ms time.Duration) {
	bits := byte(0x01 << uint(x))

	b.writeByte(y, bits)

	time.Sleep(ms * time.Millisecond)

	b.writeByte(y, 0x00)
}

// Flip flips a stone at (x, y)
func (b *MrBoard) Flip(x int, y int) {
	b.highWhile(x, y, FLIPTIME)
}
