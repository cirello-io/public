package main

import (
	"fmt"
	"os"
)

func main() {
	var cmd string

	if len(os.Args) == 1 {
		cmd = "help"
	} else {
		cmd = os.Args[1]
	}

	switch cmd {
	case "init":
		Init()

	case "add":
		Add()

	case "up":
		Up()

	case "sync":
		Sync()

	case "report":
		Report()

	case "remove":
		Remove()

	case "down":
		Down()

	default:
		help()
	}
}

func help() {
	fmt.Printf(`ezwg is a simple tool to manage a centralised wireguard VPN.

Usage: ezwg <cmd>

Available commands:

	init   : Create %[1]s containing default configuration + new keys without loading. Edit to taste.
	add    : Add a new peer + sync
	up     : Create the interface, run pre/post up, sync
	report : Generate a JSON status report to the location configured in %[1]s.
	remove : Remove a peer by hostname provided as argument + sync
	down   : Destroy the interface, run pre/post down
	sync   : Update wireguard configuration from %[1]s after validating


ezwg version %[2]s
commit %[3]s
built %[4]s

`, ConfigFile, Version, GitCommit, BuildDate)
}
