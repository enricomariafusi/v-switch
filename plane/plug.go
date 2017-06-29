package plane

import (
	"V-switch/conf"
	"V-switch/crypt"
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

	seedAddress := conf.GetConfigItem("SEED")

	if seedAddress == "MASTER" {
		log.Println("[PLANE][PLUG]: NO SEED, We are master node: Yay! ")
	} else {

		log.Println("[PLANE][PLUG]: Starting SEED to: ", seedAddress)
		go SeedingTask(seedAddress)
	}

}

func SeedingTask(remote string) {

	cycle, _ := strconv.Atoi(conf.GetConfigItem("TTL"))
	log.Println("[PLANE][PLUG] TTL is:", cycle)

	var e error

	for e == nil {
		_, e = net.ParseMAC(VSwitch.HAddr)
		log.Println("[PLANE][PLUG] Waiting 3 seconds the MAC is there")
		time.Sleep(3 * time.Second)

	}

	for VSwitch.Server == nil {

		log.Println("[PLANE][PLUG] Waiting 3 seconds the UDP server is running")
		time.Sleep(3 * time.Second)

	}

	log.Println("[PLANE][PLUG][ANNOUNCE] Our address is :", VSwitch.HAddr)

	tmpAnnounce := VSwitch.HAddr + "|" + VSwitch.IPAdd
	// create a fake announceTLV
	tmpTlv := tools.CreateTLV("A", []byte(tmpAnnounce))
	encTlv := crypt.FrameEncrypt([]byte(VSwitch.SwID), tmpTlv)

	log.Printf("[PLANE][PLUG][ANNOUNCE] Sending announce of %s to %s: [%s]", VSwitch.HAddr, remote, tmpAnnounce)
	DispatchUDP(encTlv, remote)

	tmpTlv = tools.CreateTLV("Q", []byte(VSwitch.HAddr))
	encTlv = crypt.FrameEncrypt([]byte(VSwitch.SwID), tmpTlv)
	log.Printf("[PLANE][PLUG][ANNOUNCE] Query %s for addresses: done", remote)
	DispatchUDP(encTlv, remote)

	for {

		// announces everybody + self to everybody
		for alienmac := range VSwitch.SPlane {
			AnnounceLocal(alienmac)
			for destmac := range VSwitch.SPlane {
				if alienmac != destmac {
					AnnounceAlien(alienmac, destmac)
				}
			}

		}

		time.Sleep(time.Duration(cycle) * time.Second)

	}

}
