package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

func main() {
	gbot := gobot.NewGobot()

	e := edison.NewEdisonAdaptor("edison")
	pin1 := gpio.NewDirectPinDriver(e, "pin1", "13")
	pin2 := gpio.NewDirectPinDriver(e, "pin2", "12")

	state := true

	work := func() {
		gobot.Every(1000*time.Millisecond, func() {
			if state {
				pin2.Off()
				pin1.On()
			} else {
				pin1.Off()
				pin2.On()
			}
			state = !state
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{e},
		[]gobot.Device{pin1, pin2},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
