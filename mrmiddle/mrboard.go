package mrmiddle

import (
	"fmt"
	"log"
	"time"

	"gobot.io/x/gobot/platforms/intel-iot/edison"
)

// MrMiddle is Magic Reversi's middle ware object
type MrMiddle struct {
	e *edison.Adaptor
}

// NewMrMiddle returns MrMiddle instance
func NewMrMiddle() (mm *MrMiddle, err error) {
	mm = &MrMiddle{}

	mm.e = edison.NewAdaptor()

	err = mm.e.Connect()

	if checkError(err) {
		return nil, wrapError(err)
	}

	err = mm.Init()

	if checkError(err) {
		return nil, wrapError(err)
	}

	return
}

type row [8]bool

func (r row) reversed() (reversed row) {
	for i, v := range r {
		reversed[len(reversed)-i-1] = v
	}

	return
}

func checkError(err error) bool {
	return err != nil
}

func wrapError(err error) error {
	return fmt.Errorf("Middleware Error: %s", err)
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
func (mm *MrMiddle) Init() (err error) {
	log.Println("initialize circuit...")

	for _, addr := range EXIA {
		mm.e.I2cStart(addr)

		//　Initialize IOCON
		err = mm.e.I2cWrite(addr, []byte{IOCON, 0x00})

		if checkError(err) {
			return wrapError(err)
		}

		// Initialize IODIR as read
		err = mm.e.I2cWrite(addr, []byte{IODIRA, 0x00})

		if checkError(err) {
			return wrapError(err)
		}
		err = mm.e.I2cWrite(addr, []byte{IODIRB, 0x00})

		if checkError(err) {
			return wrapError(err)
		}
	}

	for _, addr := range EXOA {
		mm.e.I2cStart(addr)

		//　Initialize IOCON
		err = mm.e.I2cWrite(addr, []byte{IOCON, 0x00})

		if checkError(err) {
			return wrapError(err)
		}

		// Initialize IODIR as read
		err = mm.e.I2cWrite(addr, []byte{IODIRA, 0xFF})

		if checkError(err) {
			return wrapError(err)
		}
		err = mm.e.I2cWrite(addr, []byte{IODIRB, 0xFF})

		if checkError(err) {
			return wrapError(err)
		}
	}

	mm.writeAllLow()

	return
}

// read given y line
func (mm *MrMiddle) readLine(y int) (r row, err error) {
	addr, gpio := y2AddrAndGpio(y)

	mm.e.I2cStart(addr)

	err = mm.e.I2cWrite(addr, []byte{byte(gpio)})

	if checkError(err) {
		return row{}, wrapError(err)
	}

	data, err := mm.e.I2cRead(addr, 1)

	if checkError(err) {
		return row{}, wrapError(err)
	}

	r = byte2Row(data[0])

	if y%2 == 0 {
		r = r.reversed()
	}

	return
}

// read both A and B bits of given address of the Expander
func (mm MrMiddle) readAB(addr int) (byteSet [2]row, err error) {
	mm.e.I2cStart(addr)

	err = mm.e.I2cWrite(addr, []byte{GPIOA})

	if checkError(err) {
		return [2]row{}, wrapError(err)
	}

	data, err := mm.e.I2cRead(addr, 2)

	if checkError(err) {
		return [2]row{}, wrapError(err)
	}

	byteSet[0] = byte2Row(data[0]).reversed()
	byteSet[1] = byte2Row(data[1])

	return
}

// ReadWholeBoard returns whole board status
func (mm *MrMiddle) readWholeBoard() (byteSet [8]row, err error) {
	for i, addr := range EXIA {
		data, err := mm.readAB(addr)

		if checkError(err) {
			return [8]row{}, wrapError(err)
		}

		byteSet[2*i] = data[0].reversed()
		byteSet[2*i+1] = data[1]
	}

	return
}

// GetInput waits until board state changes and return x, y
func (mm *MrMiddle) GetInput() (int, int, error) {
	old, err := mm.readWholeBoard()

	if checkError(err) {
		return -1, -1, wrapError(err)
	}

	for {
		crr, err := mm.readWholeBoard()

		if checkError(err) {
			return -1, -1, wrapError(err)
		}

		if crr != old {
			for i, r := range crr {
				for j, v := range r {
					// if current status is True and old status is False then reutrn (x, y)
					if v && !old[i][j] {
						return j + 1, i + 1, nil
					}
				}
			}
		}

		time.Sleep(POLLTIME)
	}
}

// write byte data to designated line
func (mm *MrMiddle) writeByte(y int, v byte) (err error) {
	addr, gpio := y2AddrAndGpio(y - 1)

	err = mm.e.I2cStart(addr)

	if checkError(err) {
		return wrapError(err)
	}

	data := []byte{byte(gpio), v}
	err = mm.e.I2cWrite(addr, data)

	if checkError(err) {
		return wrapError(err)
	}

	return
}

func (mm *MrMiddle) writeAllLow() (err error) {
	for y := 0; y < 8; y++ {
		err = mm.writeByte(y, 0x00)

		if checkError(err) {
			return wrapError(err)
		}
	}

	return
}

// HighWhile make (x, y) to High while ms[msec]
func (mm *MrMiddle) highWhile(x int, y int, ms time.Duration) (err error) {
	bits := byte(0x01 << uint(x-1))

	err = mm.writeByte(y, bits)

	if checkError(err) {
		return wrapError(err)
	}

	time.Sleep(ms * time.Millisecond)

	err = mm.writeByte(y, 0x00)

	if checkError(err) {
		return wrapError(err)
	}

	return
}

// Flip flips a stone at (x, y)
func (mm *MrMiddle) Flip(x int, y int, pd Pole) (err error) {
	err = mm.highWhile(x, y, FLIPTIME)

	if checkError(err) {
		return wrapError(err)
	}

	return
}
