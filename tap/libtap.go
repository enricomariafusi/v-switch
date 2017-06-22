package tap

import (
	"V-switch/conf"
	"V-switch/plane"
	"V-switch/tools"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type Vswitchdevice struct {
	devicename string
	mtu        int
	frame      []byte
	Realif     *TapConn
	err        error
	mac        string
}

//This will represent the tap device when exported.
var VDev Vswitchdevice

func init() {

	VDev.SetDeviceConf()
	go VDev.ReadFrameThread()  //this is blocking so it must be a new thread
	go VDev.WriteFrameThread() //thread which writes frames into the interface
}

func (vd *Vswitchdevice) SetDeviceConf() {

	if vd.mtu, vd.err = strconv.Atoi(conf.GetConfigItem("MTU")); vd.err != nil {
		log.Printf("[TAP] Cannot get MTU from conf: <%s>", vd.err)
		vd.mtu = 1500
		vd.frame = make([]byte, vd.mtu+14)
		log.Printf("[TAP] Using the default of 1500. Hope is fine.")
	} else {
		vd.frame = make([]byte, vd.mtu+14)
		log.Printf("[TAP] MTU SET TO: %v", vd.mtu)
	}

	vd.devicename = conf.GetConfigItem("DEVICENAME")
	log.Printf("[TAP] Devicename in conf is: %v", vd.devicename)

	//	vd.Realif, vd.err = newTAP(vd.devicename)

	vd.Realif = new(TapConn)
	vd.err = vd.Realif.Open(uint(vd.mtu), vd.devicename)
	if vd.err != nil {
		log.Printf("[TAP][ERROR] Error creating tap: <%s>", vd.err)
		log.Println("[TAP][ERROR] Are you ROOT?")
		os.Exit(1)

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

		cleanframe := tools.CleanFrame(vd.frame)
		log.Printf("[TAP][READ] I/O Size: %d , Raw size %d, Clean size %d", n, len(vd.frame), len(cleanframe))
		plane.TapToPlane <- cleanframe
		log.Printf("[TAP][READ] Frame sent to Plane")

	}

}

func (vd *Vswitchdevice) WriteFrameThread() {

	log.Printf("[TAP][WRITE] Tap writing thread started")

	for n_frame := range plane.PlaneToTap {

		n, err := vd.Realif.Write(n_frame)
		if err != nil {
			log.Printf("[TAP][WRITE][ERROR] Error writing %d bytes to %s : %s", len(n_frame), vd.devicename, err.Error())
		} else {
			log.Printf("[TAP][WRITE] %d long frame of %d , from %s -> %s  to dev %s", n, len(n_frame), tools.MACSource(n_frame).String(), tools.MACDestination(n_frame).String(), vd.devicename)
		}

	}

}

//EngineStart triggers the init function in the package tap
func EngineStart() {

	log.Println("[TAP] Tap Engine Init")

}
