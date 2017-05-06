package tap

import (
	"V-switch/conf"
	"fmt"
	"log"
	"strconv"

	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
)

type vswitchdevice struct {
	devicename string
	mtu        int
	deviceif   water.Config
	frame      ethernet.Frame
	realif     *water.Interface
	err        error
}

//This will represent the tap device when exported.
var VDev vswitchdevice

func init() {

	VDev.SetDeviceConf()
	VDev.tapDeviceInit()
}

func (vd *vswitchdevice) SetDeviceConf() {

	if vd.mtu, vd.err = strconv.Atoi(conf.VConfig["MTU"]); vd.err != nil {
		log.Printf("[TAP] Cannot get MTU from conf: <%s>", vd.err)
	} else {
		vd.frame.Resize(vd.mtu)
		log.Printf("[TAP] MTU SET TO: %v", vd.mtu)
	}

	vd.devicename = conf.VConfig["DEVICENAME"]
	log.Printf("[TAP] Devicename in conf is: %v", vd.devicename)

	vd.deviceif = water.Config{
		DeviceType: water.TAP,
	}

	vd.deviceif.Name = vd.devicename

}

//creates a TAP device with name specified as argument
// just do ;
//sudo ip addr add 10.1.0.10/24 dev <tapname>
//sudo ip link set dev <tapname> up
//ping -c1 -b 10.1.0.255
func (vd *vswitchdevice) tapDeviceInit() {

	defer func() {
		if e := recover(); e != nil {
			log.Println("[TAP][EXCEPTION] OH, SHIT.")
			err, ok := e.(error)
			if !ok {
				err = fmt.Errorf("[TAPDRV]: %v", e)
			}
			log.Printf("[TAP][EXCEPTION] Error: <%s>", err)

		}
	}()

	vd.realif, vd.err = water.New(vd.deviceif)
	if vd.err != nil {
		log.Printf("[TAP][ERROR] Error creating tap: <%s>", vd.err)
		log.Println("[TAP][ERROR] Are you ROOT?")
	} else {
		log.Printf("[TAP] Success creating tap: <%s>", vd.devicename)
	}

	for {
		var n int
		n, vd.err = vd.realif.Read([]byte(vd.frame))
		if vd.err != nil {
			log.Printf("[TAP] Error reading tap: <%s>", vd.err)
		} else {
			vd.frame = vd.frame[:n]
			log.Printf("Dst: %s\n", vd.frame.Destination())
			log.Printf("Src: %s\n", vd.frame.Source())
			log.Printf("Ethertype: % x\n", vd.frame.Ethertype())
			log.Printf("Payload: % x\n", vd.frame.Payload())
		}
	}

}

//EngineStart triggers the init function in the package tap
func EngineStart() {

	log.Println("[TAP] Tap Engine Starting...")

}
