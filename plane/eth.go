package plane

import (
	"V-switch/conf"
	"V-switch/crypt"

	"V-switch/tools"
	"log"
	"net"
	"strings"
	"time"
)

func init() {

	go TapInterpreterThread()

}

func TapInterpreterThread() {

	log.Printf("[PLANE][ETH] Ethernet Thread initialized")

	for {

		_, e := net.ParseMAC(VSwitch.HAddr)

		if e != nil {
			log.Println("[PLANE][ETH] Waiting 10 seconds device is there")
			time.Sleep(10 * time.Second)

		} else {
			break
		}

	}

	var myframe []byte
	var mymacaddr string
	var mytlv []byte
	var encframe []byte
	var ekey []byte

	for {

		myframe = <-TapToPlane
		log.Printf("[PLANE][ETH] Read %d Bytes frame from channel", len(myframe))
		mymacaddr = tools.MACDestination(myframe).String()
		ekey = []byte(conf.GetConfigItem("SWITCHID"))
		encframe = crypt.FrameEncrypt(ekey, myframe)
		mytlv = tools.CreateTLV("F", encframe)

		if tools.IsMacBcast(mymacaddr) {

			for mac, _ := range VSwitch.SPlane {

				DispatchTLV(mytlv, strings.ToUpper(mac))
			}

		} else {

			DispatchTLV(mytlv, strings.ToUpper(mymacaddr))

		}

	}

}
