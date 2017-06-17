package plane

import (
	"V-switch/conf"
	"V-switch/crypt"
	"V-switch/tools"
	"log"

	"github.com/songgao/packets/ethernet"
)

func init() {

	go TapInterpreterThread()
	log.Printf("[PLANE][ETH] Ethernet Thread initialized")
}

func TapInterpreterThread() {

	var myframe ethernet.Frame
	var mymacaddr string
	var mytlv []byte
	var encframe []byte
	var ekey []byte

	for {

		myframe = <-TapToPlane
		mymacaddr = myframe.Destination().String()
		ekey = []byte(conf.GetConfigItem("SWITCHID"))
		encframe = crypt.FrameEncrypt(ekey, myframe)
		mytlv = tools.CreateTLV("F", encframe)

		if tools.IsMacBcast(mymacaddr) {

			for mac, _ := range VSwitch.Ports {

				DispatchTLV(mytlv, mac)
			}

		} else {

			DispatchTLV(mytlv, mymacaddr)

		}

	}

}
