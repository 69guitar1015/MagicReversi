package main

import (
	"log"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

const (
	// IODIRA is 0x00
	IODIRA = 0x00 + iota
	// IODIRB is 0x01
	IODIRB
	// IPOLA is 0x02
	IPOLA
	// IPOLB is 0x03
	IPOLB
	// GPINTENA is 0x04
	GPINTENA
	// GPINTENB is 0x05
	GPINTENB
	// DEFVALA is 0x06
	DEFVALA
	// DEFVALB is 0x07
	DEFVALB
	// INTCONA is 0x08
	INTCONA
	// INTCONB is 0x09
	INTCONB
	// IOCON is 0x0A
	IOCON
	// IOCON2 is 0x0B
	IOCON2
	// GPPUA is 0x0C
	GPPUA
	// GPPUB is 0x0D
	GPPUB
	// INTFA is 0x0E
	INTFA
	// INTFB is 0x0F
	INTFB
	// INTCAPA is 0x10
	INTCAPA
	// INTCAPB is 0x11
	INTCAPB
	// GPIOA is 0x12
	GPIOA
	// GPIOB is 0x13
	GPIOB
	// OLATA is 0x14
	OLATA
	// OLATB is 0x15
	OLATB
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Edison Adapter
	e := edison.NewEdisonAdaptor("edison")

	// I/O Expander(MCP23017) Drivers
	// 4 bits of the head address is fixed as 0100
	// bottom 3 bits is modifiable by setting A0, A1, A2 pin of the device
	// hole address becomes |0|1|0|0|A2|A1|A0|
	// we use this as A2=0 devices are for read and A2=1 devices are for write

	work := func() {
		var data []byte
		var err error

		err = e.I2cStart(0x20)
		checkError(err)

		defer e.Finalize()

		log.Println("initialize...")
		err = e.I2cWrite(0x20, []byte{IODIRA, 0x00})
		checkError(err)
		err = e.I2cWrite(0x20, []byte{IODIRB, 0xFF})
		checkError(err)
		err = e.I2cWrite(0x20, []byte{IOCON, 0x00})
		checkError(err)

		log.Println("start reading...")
		for {
			err = e.I2cWrite(0x20, []byte{GPIOA})
			checkError(err)

			data, err = e.I2cRead(0x20, 2)
			checkError(err)

			log.Printf("data : Aが0b%08b と　Bが0b%08b\n", data[0], data[1])

			time.Sleep(1000 * time.Millisecond)
		}
	}

	robot := gobot.NewRobot("MagicReversi",
		[]gobot.Connection{e},
		[]gobot.Device{},
		work,
	)

	gbot := gobot.NewGobot()

	gbot.AddRobot(robot)

	gbot.Start()
}
