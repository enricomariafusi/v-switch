package plane

import (
	"V-switch/conf"
	"V-switch/tools"
	"log"
	"net"
	"strings"
)

type Sport struct {
	EndPoint net.UDPAddr
	Socket   net.UDPConn
	EthIP    net.IPAddr
}

type vswitchplane struct {
	//ports map[net.HardwareAddr.String()]net.UDPAddr.String()
	//Mapping from MAC address to UDP address
	//Ports map[string]string
	//Conns map[string]net.Conn
	SPlane map[string]Sport
	HAddr  string
	Fqdn   string
	IPAdd  string
	SwID   string
	DevN   string
}

//V-Switch will be exported to UDP and to TAP
var VSwitch vswitchplane

func init() {

	log.Printf("[PLANE][PLANE] TABLES INITIALIZED")

	VSwitch.SPlane = make(map[string]Sport)

	log.Printf("[PLANE][PLANE] PORTS: %b", len(VSwitch.SPlane))

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

	_, exists := sw.SPlane[hwaddr]

	return exists
}

//Removes MAC from switch
func (sw *vswitchplane) RemoveMAC(mac string) {

	hwaddr := strings.ToUpper(mac)

	if sw.macIsKnown(hwaddr) {
		log.Printf("[PLANE][PORT][DELETE] [ %s ] Deleted from plane", hwaddr)
		delete(sw.SPlane, hwaddr)
	} else {
		log.Printf("[PLANE][PORT][DELETE] [ %s ] Non existing, cannot delete from plane", hwaddr)
	}

}

//Adds a new port into the plane
func (sw *vswitchplane) addPort(mac string, endpoint net.UDPAddr) {

	hwaddr := strings.ToUpper(mac)

	_, err := net.ParseMAC(hwaddr)
	if err != nil {
		log.Printf("[PLANE][PORT][ERROR] [ %s ] is not a valid MAC address: %s", hwaddr, err.Error())
		return
	}

	if hwaddr == sw.HAddr {
		log.Printf("[PLANE][PORT][NOOP] [ %s ] = %s : no need to add", hwaddr, sw.HAddr)
		return
	}

	var port Sport

	if sw.macIsKnown(hwaddr) {
		port.Socket = sw.SPlane[hwaddr].Socket
		port.EthIP = sw.SPlane[hwaddr].EthIP
		sw.RemoveMAC(hwaddr)
	}

	port.EndPoint = endpoint
	sw.SPlane[hwaddr] = port
	log.Printf("[PLANE][PORT][ANN] Updated port : MAC %s -> %s ", hwaddr, endpoint.String())

}

//Adds a new remoteip into the plane
func (sw *vswitchplane) addRemoteIp(mac string, remoteip net.IPAddr) {

	hwaddr := strings.ToUpper(mac)

	_, err := net.ParseMAC(hwaddr)
	if err != nil {
		log.Printf("[PLANE][PORT][ERROR] [ %s ] is not a valid MAC address: %s", hwaddr, err.Error())
		return
	}

	if hwaddr == sw.HAddr {
		log.Printf("[PLANE][PORT][NOOP] [ %s ] = %s : no need to add", hwaddr, sw.HAddr)
		return
	}

	var port Sport

	if sw.macIsKnown(hwaddr) {
		port.Socket = sw.SPlane[hwaddr].Socket
		port.EndPoint = sw.SPlane[hwaddr].EndPoint
		sw.RemoveMAC(hwaddr)
	}

	port.EthIP = remoteip
	sw.SPlane[hwaddr] = port
	log.Printf("[PLANE][REMOTEIP][ANN] Updated port : MAC %s -> %s ", hwaddr, remoteip.String())
	tools.AddARPentry(hwaddr, remoteip.String(), sw.DevN)

}

//Adds a new conn into the plane
func (sw *vswitchplane) addConn(mac string, conn net.UDPConn) {

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

	var port Sport

	if sw.macIsKnown(mac) {
		port.EndPoint = sw.SPlane[hwaddr].EndPoint
		port.EthIP = sw.SPlane[hwaddr].EthIP

		old_socket := sw.SPlane[hwaddr].Socket
		old_socket.Close()
		sw.RemoveMAC(hwaddr)
	}

	port.Socket = conn
	sw.SPlane[hwaddr] = port
	log.Printf("[PLANE][SOCKET][ANN] Updated Connection : MAC %s -> %s ", hwaddr, conn.RemoteAddr().Network())

}

//PlaneInit is just a wrapper for starting the init
func PlaneInit() {

	log.Println("[PLANE] PLANE Engine Init")

}
