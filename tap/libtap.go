package tap

import (
	"log"

	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
)

//creates a TAP device with name specified as argument
// just do ;
//sudo ip addr add 10.1.0.10/24 dev <tapname>
//sudo ip link set dev <tapname> up
//ping -c1 -b 10.1.0.255
func tapDeviceInit(devname string) {

	config := water.Config{
		DeviceType: water.TAP,
	}
	config.Name = devname

	ifce, err := water.New(config)
	if err != nil {
		log.Println(err)
	}
	var frame ethernet.Frame

	for {
		frame.Resize(1500)
		n, err := ifce.Read([]byte(frame))
		if err != nil {
			log.Println(err)
		}
		frame = frame[:n]
		log.Printf("Dst: %s\n", frame.Destination())
		log.Printf("Src: %s\n", frame.Source())
		log.Printf("Ethertype: % x\n", frame.Ethertype())
		log.Printf("Payload: % x\n", frame.Payload())
	}

}

func TapEngineStart() {

	tapDeviceInit("switch0")

}
