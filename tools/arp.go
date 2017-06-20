package tools

import (
	"log"
	"net"
	"os/exec"
	"strings"
)

func AddARPentry(mac string, ip string, dev string) {

	hwaddr := strings.ToUpper(mac)

	_, err := net.ParseMAC(hwaddr)
	if err != nil {
		log.Printf("[TOOLS][ARP] [ %s ] is not a valid MAC address: %s", hwaddr, err.Error())
		return
	}

	if net.ParseIP(ip) == nil {
		log.Printf("[TOOLS][ARP] [ %s ] is not a valid IP address: %s", ip, err.Error())
		return
	}

	//ip neigh add 130.122.130.77 lladdr A2:77:35:FA:1E:F5 dev tap0
	// noarp because we don't need gratuitous  bullshit attacks
	arpcmd := exec.Command("ip", "neigh", "replace", ip, "lladdr", hwaddr, "dev", dev, "nud", "noarp")

	err = arpcmd.Run()

	if err != nil {
		log.Printf("[TOOLS][ARP] After executing  %q : %s", arpcmd.Args, err.Error())
	} else {
		log.Printf("[TOOLS][ARP] Executed   %q", arpcmd.Args)
	}

}
