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
		log.Printf("[TOOLS][ARP][ADD] [ %s ] is not a valid MAC address: %s", hwaddr, err.Error())
		return
	}

	if net.ParseIP(ip) == nil {
		log.Printf("[TOOLS][ARP][ADD] [ %s ] is not a valid IP address: %s", ip, err.Error())
		return
	}

	//ip neigh add 130.122.130.77 lladdr A2:77:35:FA:1E:F5 dev tap0
	// noarp because we don't need gratuitous  bullshit attacks
	arpcmd := exec.Command("ip", "neigh", "replace", ip, "lladdr", hwaddr, "dev", dev, "nud", "noarp")

	err = arpcmd.Run()

	if err != nil {
		log.Printf("[TOOLS][ARP][ADD] After executing  %q : %s", arpcmd.Args, err.Error())
	} else {
		log.Printf("[TOOLS][ARP][ADD] Executed   %q", arpcmd.Args)
	}

}

func DelARPentry(ip string, dev string) {

	if net.ParseIP(ip) == nil {
		log.Printf("[TOOLS][ARP][DEL] [ %s ] is not a valid IP address", ip)
		return
	}

	// ip  neigh del 3000::a0a:a3a dev eth1

	arpcmd := exec.Command("ip", "neigh", "del", ip, "dev", dev, "nud", "noarp")

	err := arpcmd.Run()

	if err != nil {
		log.Printf("[TOOLS][ARP][DEL] After executing  %q : %s", arpcmd.Args, err.Error())
	} else {
		log.Printf("[TOOLS][ARP][DEL] Executed   %q", arpcmd.Args)
	}

}
