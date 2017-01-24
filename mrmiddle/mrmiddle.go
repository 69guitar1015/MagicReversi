package mrmiddle

import (
	"fmt"
	"log"

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

	return
}

type row [8]bool

func (r row) reversed() (reversed row) {
	for i, v := range r {
		reversed[len(reversed)-i-1] = v
	}

	return
}

func (r row) toByte() (b byte) {
	b = 0

	for i, v := range r {
		if v {
			b += 0x01 << uint(i)
		}
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
		err = mm.e.I2cWrite(addr, []byte{IODIRA, 0xFF})

		if checkError(err) {
			return wrapError(err)
		}
		err = mm.e.I2cWrite(addr, []byte{IODIRB, 0xFF})

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
		err = mm.e.I2cWrite(addr, []byte{IODIRA, 0x00})

		if checkError(err) {
			return wrapError(err)
		}
		err = mm.e.I2cWrite(addr, []byte{IODIRB, 0x00})

		if checkError(err) {
			return wrapError(err)
		}
	}

	mm.writeAllLow()

	return
}
