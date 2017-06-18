package plane

import (
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
}

//V-Switch will be exported to UDP and to TAP
var VSwitch vswitchplane

func init() {

	log.Printf("[PLANE][PLANE] TABLES INITIALIZED")
	VSwitch.Ports = make(map[string]string)
	VSwitch.Conns = make(map[string]net.Conn)

	log.Printf("[PLANE][PLANE] PORTS: %b", len(VSwitch.Ports))
	log.Printf("[PLANE][PLANE] CONNS: %b", len(VSwitch.Conns))

}

//Returns true if the MAC is known
func (sw *vswitchplane) macIsKnown(mac string) bool {

	_, exists := sw.Ports[mac]

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

	_, err = net.ParseMAC(mac)
	if err != nil {
		log.Printf("[PLANE][PORT][ERROR] [ %s ] is not a valid MAC address: %s", mac, err.Error())
		return
	}

	if mac == sw.HAddr {
		log.Printf("[PLANE][PORT][NOOP] [ %s ] = %s : no need to add", mac, sw.HAddr)
		return
	}

	log.Printf("[PLANE][PORT][ANN] Updated port -> MAC %s to %s ", mac, ind)
	sw.Ports[hwaddr] = ind

}

//Adds a new conn into the plane
func (sw *vswitchplane) addConn(mac string, conn net.Conn) {

	hwaddr := strings.ToUpper(mac)

	_, err := net.ParseMAC(hwaddr)
	if err != nil {
		log.Printf("[PLANE][CONN][ERROR] [ %s ] is not a valid MAC address: %s", mac, err.Error())
		return
	}

	if mac == sw.HAddr {
		log.Printf("[PLANE][PORT][NOOP] %s = %s : no need to add", mac, sw.HAddr)
		return
	}

	log.Printf("[PLANE][CONN][NEW] Added New port -> MAC %s -> %s ", mac, conn.RemoteAddr().String())
	sw.Conns[hwaddr] = conn

}

//PlaneInit is just a wrapper for starting the init
func PlaneInit() {

	log.Println("[PLANE] PLANE Engine Init")

}
