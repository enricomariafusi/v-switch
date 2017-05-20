package plane

import (
	"V-switch/conf"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/songgao/packets/ethernet"
)

type udpframe struct {
	uframe ethernet.Frame //frame we received from the network
	addr   string         //Address to send/receive
}

type vswitchplane struct {
	//ports map[net.HardwareAddr.String()]net.UDPAddr.String()
	//Mapping from MAC address to UDP address
	ports map[string]string
	// from TAP to UDP
	ToBroadcast chan udpframe //Frames to be broadcasted , see tap.IsMacBcast
	ToUdpSend   chan udpframe //Frames to be sent to UDP

	//UDP to PLANE
	ToAnnounceOne   chan udpframe //Signaling to announce themselves  (MAC 13:13:13:13:13:13 )
	ToAnnounceAlien chan udpframe //Signaling to announce someone else (MAC 17:17:17:17:17:17 )
	//IEEE 802.3-2012 section 6 clause 79

	//UDP to TAP
	ToTapSend chan ethernet.Frame

	// Queue size
	queuesize int

	// Plane errors
	err error
}

//V-Switch will be exported to UDP and to TAP
var VSwitch vswitchplane

//Uframe will be exported to UDP
var UFrame udpframe

func init() {

	if VSwitch.queuesize, VSwitch.err = strconv.Atoi(conf.GetConfigItem("QUEUE")); VSwitch.err != nil {
		log.Printf("[PLANE] Cannot get QUEUE from conf: <%s>, using default 256", VSwitch.err)
		VSwitch.queuesize = 256
	}
	log.Printf("[PLANE] QUEUE SET TO: %v", VSwitch.queuesize)

	VSwitch.ports = make(map[string]string)

	VSwitch.ToBroadcast = make(chan udpframe, VSwitch.queuesize)
	VSwitch.ToUdpSend = make(chan udpframe, VSwitch.queuesize)

	VSwitch.ToAnnounceOne = make(chan udpframe, VSwitch.queuesize)
	VSwitch.ToAnnounceAlien = make(chan udpframe, VSwitch.queuesize)

	VSwitch.ToTapSend = make(chan ethernet.Frame, VSwitch.queuesize)

	log.Printf("[PLANE] QUEUES INITIALIZED")
	log.Printf("[PLANE] PORTS: %b", len(VSwitch.ports))

	go VSwitch.broadcastFrame()

}

//broadcastFrame takes from ToBroadcast
//and sends the same frame to  ToUdpSend , foreach port.
func (sw *vswitchplane) broadcastFrame() {

	log.Printf("[PLANE] Broadcast engine started")

	for {

		myudpframe := <-sw.ToBroadcast

		for _, tmp_addr := range sw.ports {
			myudpframe.addr = tmp_addr
			sw.ToUdpSend <- myudpframe
		}

	}

	log.Printf("[PLANE] Broadcast engine ends")

}

//UdpUniAnnounce is taking a frame from ToAnnounceOne
//if the MAC is known, it will update the port map
//if the MAC is new, it will ALSO put the frame into ToBroadcast
func (sw *vswitchplane) udpUniAnnounce() {

}

//UdpAlienAnnounce is taking a frame from ToAnnounceAlien
//takes the address from the payload instead of the senderaddress
//updates/adds the ports,
//and puts  frame into the ToBroadcast queue
func (sw *vswitchplane) udpAlienAnnounce() {

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

//PlaneInit is just a wrapper for starting the init
func PlaneInit() {

	log.Println("[PLANE] PLANE Engine Init")

}
