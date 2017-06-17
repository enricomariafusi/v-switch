package plane

import (
	"V-switch/conf"
	"V-switch/crypt"
	"V-switch/tools"
	"log"
	"net"
	"strings"
)

func init() {

	go TLVInterpreter()

}

func TLVInterpreter() {

	var my_tlv []byte
	log.Println("[PLANE][TLV] TLV receive thread starts")

	for {

		my_tlv = <-UdpToPlane

		typ, _, payload := tools.UnPackTLV(my_tlv)

		switch typ {

		// it is a frame
		case "F":
			PlaneToTap <- crypt.FrameDecrypt([]byte(conf.GetConfigItem("SWITCHID")), payload)
			// someone is announging itself
		case "A":
			announce := crypt.FrameDecrypt([]byte(conf.GetConfigItem("SWITCHID")), payload)
			fields := strings.Split(string(announce), "|")
			if len(fields) == 2 {
				VSwitch.addPort(fields[0], fields[1])
				UDPCreateConn(fields[0], fields[1])
			}
		case "Q":
			sourcemac := crypt.FrameDecrypt([]byte(conf.GetConfigItem("SWITCHID")), payload)
			for alienmac, _ := range VSwitch.Ports {
				AnnounceAlien(alienmac, string(sourcemac))
			}

		default:
			log.Println("[PLANE][TLV][INTERPRETER] Unknown type, discarded: ", typ)

		}

	}

}

func UDPCreateConn(mac string, remote string) {

	_, open_already := VSwitch.Conns[mac]

	if open_already {
		return
	}

	log.Println("[PLANE][TLV]: Creating port with: ", remote)

	ServerAddr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		log.Println("[PLANE][TLV] Bad destination address ", remote, ":", err.Error())
		return
	}

	LocalAddr, err := net.ResolveUDPAddr("udp", tools.GetLocalIp()+":0")
	if err != nil {
		log.Println("[PLANE][TLV] Cannot find local port to bind ", remote, ":", err.Error())
		return
	}

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)

	if err != nil {
		log.Println("[PLANE][TLV] Error connecting with ", remote, ":", err.Error())
		return
	}
	log.Println("[PLANE][TLV] Success connecting with ", remote)

	VSwitch.addConn(mac, Conn)

	AnnounceLocal(mac)

}

func DispatchTLV(mytlv []byte, mac string) {

	_, open_already := VSwitch.Conns[mac]

	if open_already {

		go VSwitch.Conns[mac].Write([]byte(mytlv))

	} else {
		log.Println("[PLANE][TLV][DISPATCH] cannot dispatch, no connection available for ", mac)
		return
	}

}

func AnnounceLocal(mac string) {

	var myfqdn string

	if conf.ConfigItemExists("PUBLIC") {
		myfqdn = conf.GetConfigItem("PUBLIC")
	} else {
		myfqdn = tools.GetFQDN() + ":" + conf.GetConfigItem("PORT")
		log.Println("[PLANE][TLV][ANNOUNCE] dynamic hostid set to", myfqdn)
	}

	myannounce := VSwitch.HAddr + "|" + myfqdn
	mykey := conf.GetConfigItem("SWITCHID")

	myannounce_enc := crypt.FrameEncrypt([]byte(mykey), []byte(myannounce))

	tlv := tools.CreateTLV("A", myannounce_enc)

	DispatchTLV(tlv, mac)

}

func AnnounceAlien(alien_mac string, mac string) {

	var myfqdn string

	if VSwitch.macIsKnown(alien_mac) {
		myfqdn = VSwitch.Ports[alien_mac]
	} else {
		log.Println("[PLANE][TLV][ANNOUNCE][ALIEN] cannot announce unknown mac: ", alien_mac)
		return
	}

	myannounce := strings.ToUpper(alien_mac) + "|" + myfqdn
	mykey := conf.GetConfigItem("SWITCHID")

	myannounce_enc := crypt.FrameEncrypt([]byte(mykey), []byte(myannounce))

	tlv := tools.CreateTLV("A", myannounce_enc)

	DispatchTLV(tlv, mac)

}

func SendQueryToMac(mac string) {

	myannounce := VSwitch.HAddr
	mykey := conf.GetConfigItem("SWITCHID")

	myannounce_enc := crypt.FrameEncrypt([]byte(mykey), []byte(myannounce))

	tlv := tools.CreateTLV("Q", myannounce_enc)

	DispatchTLV(tlv, mac)

}
