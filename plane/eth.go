package plane

import (
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

	var e error = nil

	for e == nil {
		_, e = net.ParseMAC(VSwitch.HAddr)
		log.Println("[PLANE][ETH] Waiting 3 seconds the MAC is there")
		time.Sleep(3 * time.Second)

	}

	var myframe []byte
	var mymacaddr string
	var mytlv []byte
	var encframe []byte
	var ekey []byte

	for {

		myframe = <-TapToPlane
		log.Printf("[PLANE][ETH] Read %d Bytes frame from QUEUE TapToPlane", len(myframe))
		mymacaddr = tools.MACDestination(myframe).String()
		ekey = []byte(VSwitch.SwID)
		mytlv = tools.CreateTLV("F", myframe)
		encframe = crypt.FrameEncrypt(ekey, mytlv)

		if tools.IsMacBcast(mymacaddr) {

			for mac, _ := range VSwitch.SPlane {

				DispatchTLV(encframe, strings.ToUpper(mac))
			}

		} else {

			DispatchTLV(encframe, strings.ToUpper(mymacaddr))

		}

	}

}
