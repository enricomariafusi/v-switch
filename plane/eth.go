package plane

import (
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

	for {

		myframe = <-TapToPlane
		mymacaddr = myframe.Destination().String()

		if tools.IsMacBcast(mymacaddr) {

			mytlv := tools.CreateTLV("F", myframe)

			for mac, _ := range VSwitch.Ports {

				DispatchTLV(mytlv, mac)
			}

		} else {

			mytlv = tools.CreateTLV("F", myframe)
			DispatchTLV(mytlv, mymacaddr)

		}

	}

}
