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

func checkError(err error) {
	if err != nil {
		// log.Fatal(err)
		log.Print("Error : ")
		log.Println(err)
	}
}

// convert from byte object to boolean array
func byte2BoolArray(b byte) (arr [8]bool) {
	for i := 0; i < 8; i++ {
		arr[i] = (b&0x01 == 0x01)
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
func (b *MrBoard) readLine(y int) (bits [8]bool) {
	addr, gpio := y2AddrAndGpio(y)

	b.e.I2cStart(addr)

	err := b.e.I2cWrite(addr, []byte{byte(gpio)})
	checkError(err)

	data, err := b.e.I2cRead(addr, 1)
	checkError(err)

	bits = byte2BoolArray(data[0])

	return
}

// read both A and B bits of given address of the Expander
func (b MrBoard) readAB(addr int) (bitsSet [2][8]bool) {
	b.e.I2cStart(addr)

	err := b.e.I2cWrite(addr, []byte{GPIOA})
	checkError(err)

	data, err := b.e.I2cRead(addr, 2)
	checkError(err)

	bitsSet[0] = byte2BoolArray(data[0])
	bitsSet[1] = byte2BoolArray(data[1])

	return
}

// ReadWholeBoard returns whole board status
func (b *MrBoard) ReadWholeBoard() (bitsSet [8][8]bool) {
	for i, addr := range EXIA {
		data := b.readAB(addr)

		bitsSet[2*i] = data[0]
		bitsSet[2*i+1] = data[1]
	}

	return
}

// write byte data to designated line
func (b *MrBoard) writeBits(y int, bits byte) {
	addr, gpio := y2AddrAndGpio(y)

	b.e.I2cStart(addr)

	data := []byte{byte(gpio), bits}
	b.e.I2cWrite(addr, data)
}

func (b *MrBoard) writeAllLow() {
	for y := 0; y < 8; y++ {
		b.writeBits(y, 0x00)
	}
}

// HighWhile make (x, y) to High while ms[msec]
func (b *MrBoard) HighWhile(x int, y int, ms time.Duration) {
	bits := byte(0x01 << uint(x))

	b.writeBits(y, bits)

	time.Sleep(ms * time.Millisecond)

	b.writeBits(y, 0x00)
}
