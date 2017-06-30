package tap

import (
	"V-switch/conf"
	"V-switch/plane"
	"log"
	"net"
	"os/exec"
	"strings"

	"time"
)

func init() {

	go SetIpAddress()
}

func SetIpAddress() {

	for {

		_, e := net.ParseMAC(plane.VSwitch.HAddr)

		if e != nil {
			log.Println("[TAP][IP] Waiting 10 seconds device is there")
			time.Sleep(10 * time.Second)

		} else {
			break
		}

	}

	my_ipnetmask := plane.VSwitch.IPAdd + "/" + conf.GetConfigItem("DEVICEMASK")

	// ifcnfg := exec.Command("ifconfig", my_dev_name, plane.VSwitch.IPAdd, "netmask", my_ipnetmask, "mtu", strconv.Itoa(eth_mtu))
	ifcnfg := exec.Command("ip", "address", "add", my_ipnetmask, "dev", plane.VSwitch.DevN)

	err := ifcnfg.Run()
	if err != nil {
		log.Printf("[TAP][IP] Error executing  %s: %s", strings.Join(ifcnfg.Args, " "), err.Error())
	} else {
		log.Printf("[TAP][IP] Executed   %s", strings.Join(ifcnfg.Args, " "))
	}

}
