package plane

import (
	"V-switch/conf"
	"V-switch/tools"
	"log"
	"net"
	"strings"
)

type vswitchplane struct {
	//ports map[net.HardwareAddr.String()]net.UDPAddr.String()
	//Mapping from MAC address to UDP address
	Ports map[string]string
	Conns map[string]net.Conn
	HAddr string
	Fqdn  string
	IPAdd string
	SwID  string
	DevN  string
}

//V-Switch will be exported to UDP and to TAP
var VSwitch vswitchplane

func init() {

	log.Printf("[PLANE][PLANE] TABLES INITIALIZED")
	VSwitch.Ports = make(map[string]string)
	VSwitch.Conns = make(map[string]net.Conn)

	log.Printf("[PLANE][PLANE] PORTS: %b", len(VSwitch.Ports))
	log.Printf("[PLANE][PLANE] CONNS: %b", len(VSwitch.Conns))

	if conf.ConfigItemExists("PUBLIC") {
		VSwitch.Fqdn = conf.GetConfigItem("PUBLIC")
		log.Println("[PLANE] dynamic hostid set to", VSwitch.Fqdn)
	} else {
		VSwitch.Fqdn = tools.GetFQDN() + ":" + conf.GetConfigItem("PORT")
		log.Println("[PLANE] dynamic hostid set to", VSwitch.Fqdn)
	}

	VSwitch.IPAdd = conf.GetConfigItem("DEVICEADDR")
	VSwitch.DevN = conf.GetConfigItem("DEVICENAME")
	VSwitch.SwID = conf.GetConfigItem("SWITCHID")

}

//Returns true if the MAC is known
func (sw *vswitchplane) macIsKnown(mac string) bool {

	hwaddr := strings.ToUpper(mac)

	_, exists := sw.Ports[hwaddr]

	return exists
}

//Adds a new port into the plane
func (sw *vswitchplane) addPort(mac string, ind string) {

	hwaddr := strings.ToUpper(mac)

	_, err := net.ResolveUDPAddr("udp", ind)
	if err != nil {
		log.Printf("[PLANE][PORT][ERROR] [ %s ] is not a valid UDP address: %s", ind, err.Error())
		return
	}

	_, err = net.ParseMAC(hwaddr)
	if err != nil {
		log.Printf("[PLANE][PORT][ERROR] [ %s ] is not a valid MAC address: %s", hwaddr, err.Error())
		return
	}

	if hwaddr == sw.HAddr {
		log.Printf("[PLANE][PORT][NOOP] [ %s ] = %s : no need to add", hwaddr, sw.HAddr)
		return
	}

	log.Printf("[PLANE][PORT][ANN] Updated port -> MAC %s to %s ", hwaddr, ind)
	sw.Ports[hwaddr] = ind

}

//Adds a new conn into the plane
func (sw *vswitchplane) addConn(mac string, conn net.Conn) {

	hwaddr := strings.ToUpper(mac)

	_, err := net.ParseMAC(hwaddr)
	if err != nil {
		log.Printf("[PLANE][CONN][ERROR] [ %s ] is not a valid MAC address: %s", hwaddr, err.Error())
		return
	}

	if mac == sw.HAddr {
		log.Printf("[PLANE][PORT][NOOP] %s = %s : no need to add", hwaddr, sw.HAddr)
		return
	}

	log.Printf("[PLANE][CONN][NEW] Added New port -> MAC %s -> %s ", hwaddr, conn.RemoteAddr().String())
	sw.Conns[hwaddr] = conn

}

//PlaneInit is just a wrapper for starting the init
func PlaneInit() {

	log.Println("[PLANE] PLANE Engine Init")

}
