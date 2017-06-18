package plane

import (
	"V-switch/conf"
	"V-switch/crypt"
	"V-switch/tools"
	"log"
	"strings"
)

func init() {

	go TapInterpreterThread()
	log.Printf("[PLANE][ETH] Ethernet Thread initialized")
}

func TapInterpreterThread() {

	var myframe []byte
	var mymacaddr string
	var mytlv []byte
	var encframe []byte
	var ekey []byte

	for {

		myframe = <-TapToPlane
		log.Printf("[PLANE][ETH] Read %d frame from channel", len(myframe))
		mymacaddr = tools.MACDestination(myframe).String()
		ekey = []byte(conf.GetConfigItem("SWITCHID"))
		encframe = crypt.FrameEncrypt(ekey, myframe)
		mytlv = tools.CreateTLV("F", encframe)

		if tools.IsMacBcast(mymacaddr) {

			for mac, _ := range VSwitch.Ports {

				DispatchTLV(mytlv, strings.ToUpper(mac))
			}

		} else {

			DispatchTLV(mytlv, strings.ToUpper(mymacaddr))

		}

	}

}
