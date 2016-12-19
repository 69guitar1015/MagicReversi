package main

import (
	"log"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
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
		err = e.I2cWrite(0x20, []byte{0x00, 0xFF})
		checkError(err)
		err = e.I2cWrite(0x20, []byte{0x01, 0xFF})
		checkError(err)

		log.Println("start reading...")
		for {
			err = e.I2cWrite(0x20, []byte{0x12})
			checkError(err)

			data, err = e.I2cRead(0x20, 2)
			checkError(err)

			log.Printf("data : %d と　%d\n", data[0], data[1])

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
