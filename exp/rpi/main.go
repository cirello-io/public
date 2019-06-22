package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	yellowLed := gpio.NewLedDriver(r, "8")
	redLed := gpio.NewLedDriver(r, "16")

	work := func() {
		yellowLed.On()
		redLed.Off()
		gobot.Every(200*time.Millisecond, func() {
			yellowLed.Toggle()
			redLed.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{r},
		[]gobot.Device{yellowLed, redLed},
		work,
	)

	robot.Start()
}
