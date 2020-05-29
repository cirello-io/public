package main

import (
	"time"
)

// ConfigFile indicates the default location of the configuration file
// and it could be overridden in future via env.
const ConfigFile = "/etc/ezwgconfig.json"

// these end up in the config file
const (
	DefaultInterfaceName = "wg0"
	DefaultReportFile    = "/var/lib/ezwgreport.json"
	DefaultListenPort    = 51820

	// keepalive always configured for clients. Set to a value likely to
	// stop most NATs from dropping the connection. Wireguard docs recommend 25
	// for most NATs
	Keepalive = 25 * time.Second

	// if last handshake (different from keepalive, see https://www.wireguard.com/protocol/)
	Timeout = 3 * time.Minute

	// when is a peer considered gone forever? (could remove)
	Expiry = 28 * time.Hour * 24
)

var (
	// populated with LDFLAGS, see do-release.sh
	Version   = "unknown"
	GitCommit = "unknown"
	BuildDate = "unknown"
)
