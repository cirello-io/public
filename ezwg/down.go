package ezwg

import (
	"github.com/vishvananda/netlink"
)

func Down() {
	conf := MustLoadezwgConfig()
	DelLink(conf)
}

func DelLink(conf *ezwgConfig) {
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = conf.InterfaceName

	link := &netlink.GenericLink{
		LinkAttrs: linkAttrs,
	}

	err := netlink.LinkDel(link)
	if err != nil {
		ExitFail("Could not delete interface '%s' (%v)", conf.InterfaceName, err)
	}
}
