package main

import (
	"math/rand"

	"github.com/sirupsen/logrus"
)

func main() {
	for {
		switch rand.Intn(3) {
		case 0:
			logrus.Debugln(renderDebugMessage())
		case 1:
			logrus.Infoln(renderInfoMessage())
		case 2:
			logrus.Errorln(renderErrorMessage())
		}
	}
}

func renderDebugMessage() string {
	return "some message"
}
func renderInfoMessage() string {
	return "some message"
}
func renderErrorMessage() string {
	return "some message"
}
