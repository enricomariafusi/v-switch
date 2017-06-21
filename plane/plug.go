package plane

import (
	"V-switch/conf"
	"V-switch/tools"
	"log"
	"net"
	"strconv"
	"time"
)

func init() {

	if conf.ConfigItemExists("SEED") == false {
		conf.SetConfigItem("SEED", "MASTER")
	}

	seed_address := conf.GetConfigItem("SEED")

	if seed_address == "MASTER" {
		log.Println("[PLANE][PLUG]: NO SEED, We are master node: Yay! ")
	} else {

		log.Println("[PLANE][PLUG]: Starting SEED to: ", seed_address)
		go SeedingTask(seed_address)
	}

}

func SeedingTask(remote string) {

	cycle, _ := strconv.Atoi(conf.GetConfigItem("TTL"))

	var e error = nil

	for e == nil {
		_, e = net.ParseMAC(VSwitch.HAddr)
		log.Println("[PLANE][PLUG] Waiting 3 seconds the MAC is there")
		time.Sleep(3 * time.Second)

	}

	log.Println("[PLANE][PLUG][ANNOUNCE] Our address is :", VSwitch.HAddr)

	tmp_tlv := tools.CreateTLV("A", []byte(VSwitch.HAddr+"|"+VSwitch.Fqdn+"|"+VSwitch.IPAdd))
	CustomDispatch(tmp_tlv, remote)
	log.Printf("[PLANE][PLUG][ANNOUNCE] Sent announce of %s to %s", VSwitch.HAddr, remote)
	tmp_tlv = tools.CreateTLV("Q", []byte(VSwitch.HAddr))
	CustomDispatch(tmp_tlv, remote)
	log.Printf("[PLANE][PLUG][ANNOUNCE] Query %s for addresses: done", remote)

	for {

		// announces everybody to everybody
		for alienmac, _ := range VSwitch.SPlane {
			for destmac, _ := range VSwitch.SPlane {
				if alienmac != destmac {
					AnnounceAlien(alienmac, destmac)
				}
			}

		}

		time.Sleep(time.Duration(cycle) * time.Second)

	}

}
