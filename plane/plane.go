package plane

import (
	"log"
	"net"
	"strings"
)

type vswitchplane struct {
	//ports map[net.HardwareAddr.String()]net.UDPAddr.String()
	//Mapping from MAC address to UDP address
	ports map[string]string

	conns map[string]net.Conn
}

//V-Switch will be exported to UDP and to TAP
var VSwitch vswitchplane

func init() {

	log.Printf("[PLANE] TABLES INITIALIZED")
	VSwitch.ports = make(map[string]string)
	VSwitch.conns = make(map[string]net.Conn)

	log.Printf("[PLANE] PORTS: %b", len(VSwitch.ports))
	log.Printf("[PLANE] CONNS: %b", len(VSwitch.conns))

}

//Returns true if the MAC is known
func (sw *vswitchplane) macIsKnown(mac net.HardwareAddr) bool {

	_, exists := sw.ports[mac.String()]

	return exists
}

//Adds a new port into the plane
func (sw *vswitchplane) addPort(mac string, ind string) {

	hwaddr := strings.ToUpper(mac)

	_, err := net.ResolveUDPAddr("udp", ind)
	if err != nil {
		log.Printf("[PLANE][PORT][ERROR] %s is not a valid UDP address", ind)
		return
	}

	_, err = net.ParseMAC(mac)
	if err != nil {
		log.Printf("[PLANE][PORT][ERROR] %s is not a valid MAC address", mac)
		return
	}

	log.Printf("[PLANE][PORT][NEW] Added New port -> MAC %s to %s ", mac, ind)
	sw.ports[hwaddr] = ind

}

//Adds a new conn into the plane
func (sw *vswitchplane) addConn(mac string, conn net.Conn) {

	_, err := net.ParseMAC(mac)
	if err != nil {
		log.Printf("[PLANE][CONN][ERROR] %s is not a valid MAC address", mac)
		return
	}

	log.Printf("[PLANE][CONN][NEW] Added New port -> MAC %s to %s ", mac, conn.RemoteAddr().String())
	sw.conns[mac] = conn

}

//PlaneInit is just a wrapper for starting the init
func PlaneInit() {

	log.Println("[PLANE] PLANE Engine Init")

}
