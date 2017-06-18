package tap

import (
	"V-switch/conf"
	"V-switch/plane"
	"log"
	"net"
	"os/exec"
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

	my_dev_name := conf.GetConfigItem("DEVICENAME")
	my_ipnetmask := conf.GetConfigItem("DEVICEMASK")

	ifcnfg := exec.Command("ifconfig", my_dev_name, plane.VSwitch.IPAdd, "netmask", my_ipnetmask)

	err := ifcnfg.Run()
	if err != nil {
		log.Printf("[TAP][IP] Error executing  %q", ifcnfg.Args)
	} else {
		log.Printf("[TAP][IP] Executed   %q", ifcnfg.Args)
	}

}
