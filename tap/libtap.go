package tap

import (
	"V-switch/conf"
	"V-switch/plane"
	"V-switch/tools"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Vswitchdevice struct {
	devicename string
	mtu        int
	frame      []byte
	Realif     Interface
	err        error
	mac        string
}

//This will represent the tap device when exported.
var VDev Vswitchdevice

func init() {

	VDev.SetDeviceConf()
	go VDev.ReadFrameThread() //this is blocking so it must be a new thread
	WDev := VDev
	go WDev.WriteFrameThread() //thread which writes frames into the interface
}

func (vd *Vswitchdevice) SetDeviceConf() {

	if vd.mtu, vd.err = strconv.Atoi(conf.GetConfigItem("MTU")); vd.err != nil {
		log.Printf("[TAP] Cannot get MTU from conf: <%s>", vd.err)
		vd.mtu = 1500
		vd.frame = make([]byte, vd.mtu+42)
		log.Printf("[TAP] Using the default of 1500. Hope is fine.")
	} else {
		vd.frame = make([]byte, vd.mtu+42)
		log.Printf("[TAP] MTU SET TO: %v", vd.mtu)
	}

	vd.devicename = conf.GetConfigItem("DEVICENAME")
	log.Printf("[TAP] Devicename in conf is: %v", vd.devicename)

	vd.Realif, vd.err = newTAP(vd.devicename)

	if vd.err != nil {
		log.Printf("[TAP][ERROR] Error creating tap: <%s>", vd.err)
		log.Println("[TAP][ERROR] Are you ROOT?")
	} else {
		tmp_if, _ := net.InterfaceByName(vd.devicename)
		vd.mac = strings.ToUpper(tmp_if.HardwareAddr.String())
		plane.VSwitch.HAddr = vd.mac
		log.Printf("[TAP] Success creating tap: <%s> at mac [%s] ", vd.devicename, vd.mac)
	}

}

//creates a TAP device with name specified as argument
// just do ;
//sudo ip addr add 10.1.0.10/24 dev <tapname>
//sudo ip link set dev <tapname> up
//ping -c1 -b 10.1.0.255
func (vd *Vswitchdevice) ReadFrameThread() {

	defer func() {
		if e := recover(); e != nil {
			log.Println("[TAP][EXCEPTION] OH, SHIT.")
			err, ok := e.(error)
			if !ok {
				err = fmt.Errorf("[TAP][DRV]: %v", e)
			}
			log.Printf("[TAP][EXCEPTION] Error: <%s>", err)

		}
	}()

	for {

		vd.ReadFrame()

	}

}

//returns mac address of the device we created
func (vd *Vswitchdevice) GetTapMac() string {

	macc, _ := net.ParseMAC(vd.mac)
	if macc != nil {
		log.Printf("[TAP] GetTapMac: %s", vd.mac)
		return strings.ToUpper(vd.mac)

	}

	log.Printf("[TAP] GetTapMac: mac address is empty, using default one")
	return "00:00:00:00:00:00"

}

func (vd *Vswitchdevice) ReadFrame() {

	var n int

	n, vd.err = vd.Realif.Read(vd.frame)

	if vd.err != nil {
		log.Printf("[TAP] Error reading tap: <%s>", vd.err)

	} else {

		//log.Printf("Src: %s , Broadcast :%t\n", vd.frame.Source(), tools.IsMacBcast(vd.frame.Source().String()))
		//log.Printf("Dst: %s , Broadcast :%t\n", vd.frame.Destination(), tools.IsMacBcast(vd.frame.Destination().String()))
		//log.Printf("Ethertype: % x\n", vd.frame.Ethertype())
		//log.Printf("Payload: % x\n", vd.frame.Payload())
		log.Printf("Size: %d\n", n)
		log.Printf("Frame(%d): % x\n", len(vd.frame), vd.frame)

		vd.frame = vd.frame[:n]

		plane.TapToPlane <- vd.frame

	}

}

func (vd *Vswitchdevice) WriteFrameThread() {

	var n_frame []byte

	log.Printf("[TAP][WRITE] Tap writing thread started")

	for {

		n_frame = <-plane.PlaneToTap

		n, err := vd.Realif.Write(n_frame)
		if err != nil {
			log.Printf("[TAP][WRITE][ERROR] Error writing to %s : %s", vd.devicename, err.Error())
		} else {
			log.Printf("[TAP][WRITE] %d long frame  from %s -> %s  to dev %s", n, tools.MACSource(n_frame).String(), tools.MACDestination(n_frame).String(), vd.devicename)
		}

	}

}

//EngineStart triggers the init function in the package tap
func EngineStart() {

	log.Println("[TAP] Tap Engine Init")

}
