package tap

import (
	"fmt"
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

	defer func() {
		if e := recover(); e != nil {
			log.Println("[TAP-EXCEPTION] OH, SHIT.")
			err, ok := e.(error)
			if !ok {
				err = fmt.Errorf("[TAPDRV]: %v", e)
			}
			log.Printf("[TAP-EXCEPTION] Error: <%s>", err)

		}
	}()

	config := water.Config{
		DeviceType: water.TAP,
	}
	config.Name = devname

	ifce, err := water.New(config)
	if err != nil {
		log.Printf("[TAP] Error creating tap: <%s>", err)
	}
	var frame ethernet.Frame

	for {
		frame.Resize(1500)
		n, err := ifce.Read([]byte(frame))
		if err != nil {
			log.Printf("[TAP] Error reading tap: <%s>", err)
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
