package tap

import (
	"V-switch/conf"
	"fmt"
	"log"
	"strconv"

	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
)

var VDev map[string]string

func init() {

	VDev = make(map[string]string)

	VDev["MTU"] = conf.VConfig["MTU"]
	log.Printf("[TAP-INIT] MTU: <%s>", VDev["MTU"])
	VDev["DEVICENAME"] = conf.VConfig["DEVICENAME"]
	log.Printf("[TAP-INIT] DEV: <%s>", VDev["DEVICENAME"])

}

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
	} else {
		log.Printf("[TAP] Success creating tap: <%s>", devname)
	}

	var frame ethernet.Frame

	if mtu, err := strconv.Atoi(VDev["MTU"]); err != nil {
		log.Printf("[TAP] Cannot get MTU from conf: <%s>", err)
	} else {
		frame.Resize(mtu)
		log.Printf("[TAP] MTU SET TO: %v", mtu)
	}

	for {

		n, err := ifce.Read([]byte(frame))
		if err != nil {
			log.Printf("[TAP] Error reading tap: <%s>", err)
		} else {
			frame = frame[:n]
			log.Printf("Dst: %s\n", frame.Destination())
			log.Printf("Src: %s\n", frame.Source())
			log.Printf("Ethertype: % x\n", frame.Ethertype())
			log.Printf("Payload: % x\n", frame.Payload())
		}
	}

}

func TapEngineStart() {

	tapDeviceInit(VDev["DEVICENAME"])

}
