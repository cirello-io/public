package main

import (
	"fmt"
	"log"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	button()
}

func button() {
	r := raspi.NewAdaptor()

	const redLedPort = "8"
	if _, err := r.DigitalPin(redLedPort, "out"); err != nil {
		log.Fatal("cannot reconfigure pin - redLed")
	}
	redLed := gpio.NewLedDriver(r, redLedPort)

	const buttonPort = "12"
	if _, err := r.DigitalPin(buttonPort, "in"); err != nil {
		log.Fatal("cannot reconfigure pin - button")
	}
	button := gpio.NewButtonDriver(r, buttonPort)

	work := func() {
		button.On(gpio.ButtonPush, func(s interface{}) {
			fmt.Println("button pushed")
			redLed.Toggle()
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{r},
		[]gobot.Device{button, redLed},
		work,
	)

	robot.Start()
}

func blink() {
	r := raspi.NewAdaptor()
	yellowLed := gpio.NewLedDriver(r, "8")
	redLed := gpio.NewLedDriver(r, "16")

	work := func() {
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
