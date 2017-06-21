package plane

import (
	"V-switch/conf"
	"V-switch/tools"
	"log"
	"net"
	"strings"
)

type Sport struct {
	EndPoint *net.UDPAddr // IP:PORT of the remote peer
	Socket   *net.UDPConn // Connection to the remote peer
	EthIP    *net.IPAddr  // Ip on the interface of remote peer.
}

type vswitchplane struct {
	//ports map[net.HardwareAddr.String()]net.UDPAddr.String()
	//Mapping from MAC address to UDP address
	//Ports map[string]string
	//Conns map[string]net.Conn
	SPlane map[string]Sport
	HAddr  string // Hardware address of the local tap device
	Fqdn   string // Public IP address if setup, or the local ip address
	IPAdd  string // local ip address
	SwID   string // deviceID
	DevN   string //name of the tap device
}

//V-Switch will be exported to UDP and to TAP
var VSwitch vswitchplane

func init() {

	log.Printf("[PLANE][PLANE] TABLES INITIALIZED")

	VSwitch.SPlane = make(map[string]Sport)

	log.Printf("[PLANE][PLANE] PORTS: %b", len(VSwitch.SPlane))

	if conf.GetConfigItem("PUBLIC") != "HOSTNAME" {
		VSwitch.Fqdn = conf.GetConfigItem("PUBLIC")
		log.Println("[PLANE] dynamic hostid set to", VSwitch.Fqdn)
	} else {
		VSwitch.Fqdn = tools.GetLocalIp() + ":" + conf.GetConfigItem("PORT")
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

		old_socket := sw.SPlane[hwaddr].Socket
		err := old_socket.Close()
		log.Printf("[PLANE][CONN][CLOSE] [ %t ] UDP Connection closed", err == nil)
		delete(sw.SPlane, hwaddr)
		log.Printf("[PLANE][PORT][DELETE] [ %s ] Deleted from plane", hwaddr)
	} else {
		log.Printf("[PLANE][PORT][DELETE] [ %s ] Non existing, cannot delete from plane", hwaddr)
	}

}

func (sw *vswitchplane) AddMac(mac string, endpoint string, remoteip string) {

	mac = strings.ToUpper(mac)

	_, err := net.ParseMAC(mac)
	if err != nil {
		log.Printf("[PLANE][PORT][ADD] [ %s ] is not a valid MAC address: %s", mac, err.Error())
		return
	}

	if (VSwitch.HAddr == mac) && (VSwitch.Fqdn == endpoint) {
		log.Printf("[PLANE][PORT][ADD] [ %s -> %s ] This is ourself, not adding", mac, endpoint)
		return
	}

	endpoint_net, errb := net.ResolveUDPAddr("udp", endpoint)
	if errb != nil {
		log.Printf("[PLANE][PORT][ADD] [ %s ] is not a valid UDP address: %s", endpoint, err.Error())
		return
	}

	remoteip_net, errc := net.ResolveIPAddr("ip", remoteip)
	if errc != nil {
		log.Printf("[PLANE][PORT][ADD] [ %s ] is not a valid IP address: %s", remoteip, err.Error())
		return
	}

	// if the MAC is known and data are the same, no need to change
	if sw.macIsKnown(mac) {
		tmp_endpoint := sw.SPlane[mac].EndPoint
		tmp_remoteip := sw.SPlane[mac].EthIP
		if (endpoint == tmp_endpoint.String()) && (remoteip == tmp_remoteip.String()) {
			return
		} else {
			sw.RemoveMAC(mac)
		}

	}

	var port Sport

	port.EndPoint = endpoint_net
	port.EthIP = remoteip_net

	LocalAddr, _ := net.ResolveUDPAddr("udp", tools.GetLocalIp()+":0")

	netsocket, neterr := net.DialUDP("udp", LocalAddr, endpoint_net)
	if neterr != nil {
		log.Println("[PLANE][PLANE][ADD] Error connecting with [", endpoint, "]:", neterr.Error())
		return
	} else {
		port.Socket = netsocket
	}

	sw.SPlane[mac] = port

	if mac != VSwitch.HAddr {

		tools.AddARPentry(mac, remoteip, VSwitch.DevN)

	}

}

//PlaneInit is just a wrapper for starting the init
func PlaneInit() {

	log.Println("[PLANE] PLANE Engine Init")

}
