package main

import (
	"golang.zx2c4.com/wireguard/wgctrl"
)

func Report() {
	conf := MustLoadezwgConfig()

	wg, err := wgctrl.New()
	check(err)
	defer wg.Close()

	dev, err := wg.Device(conf.InterfaceName)

	if err != nil {
		ExitFail("Could not retrieve device '%s' (%v)", conf.InterfaceName, err)
	}

	oldReport := MustLoadezwgReport()
	report := GenerateReport(dev, conf, oldReport)
	report.MustSave(conf.ReportFile)
}
