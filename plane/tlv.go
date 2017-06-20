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
	log.Println("[PLANE][TLV][INTERPRETER] Thread starts")

	for {

		my_tlv = <-UdpToPlane

		typ, _, payload := tools.UnPackTLV(my_tlv)

		switch typ {

		// it is a frame
		case "F":
			PlaneToTap <- crypt.FrameDecrypt([]byte(VSwitch.SwID), payload)
			// someone is announging itself
		case "A":
			announce := crypt.FrameDecrypt([]byte(VSwitch.SwID), payload)
			fields := strings.Split(string(announce), "|")
			if len(fields) == 3 {
				radd, _ := net.ResolveUDPAddr("udp", fields[1])
				rip, _ := net.ResolveIPAddr("ip", fields[2])
				VSwitch.addPort(fields[0], *radd)
				UDPCreateConn(fields[0], *radd)
				VSwitch.addRemoteIp(fields[0], *rip)
			}
		case "Q":
			sourcemac := crypt.FrameDecrypt([]byte(VSwitch.SwID), payload)
			for alienmac, _ := range VSwitch.SPlane {
				AnnounceAlien(alienmac, string(sourcemac))

			}

		default:
			log.Println("[PLANE][TLV][INTERPRETER] Unknown type, discarded: [ ", typ, " ]")

		}

	}

}

func UDPCreateConn(mac string, remote net.UDPAddr) {

	mac = strings.ToUpper(mac)

	log.Println("[PLANE][TLV]: Creating port with: ", remote)

	LocalAddr, err := net.ResolveUDPAddr("udp", tools.GetLocalIp()+":0")
	if err != nil {
		log.Println("[PLANE][TLV] Cannot find local port to bind ", remote, ":", err.Error())
		return
	}

	Conn, err := net.DialUDP("udp", LocalAddr, &remote)

	if err != nil {
		log.Println("[PLANE][TLV] Error connecting with ", remote.String(), ":", err.Error())
		return
	}
	log.Println("[PLANE][TLV] Success connecting with ", remote.String())

	VSwitch.addConn(mac, *Conn)

	AnnounceLocal(mac)

}

func DispatchTLV(mytlv []byte, mac string) {

	mac = strings.ToUpper(mac)

	if mac == VSwitch.HAddr {
		log.Printf("[PLANE][TLV][DISPATCH] %s is myself : no need to dispatch", mac)
		return
	}

	if VSwitch.macIsKnown(mac) {

		osocket := VSwitch.SPlane[mac].Socket

		osocket.Write([]byte(mytlv))
		log.Printf("[PLANE][TLV][DISPATCH] Dispatching to %s", mac)

	} else {
		log.Println("[PLANE][TLV][DISPATCH] cannot dispatch, no connection available for ", mac)
		return
	}

}

func AnnounceLocal(mac string) {

	mac = strings.ToUpper(mac)

	myannounce := VSwitch.HAddr + "|" + VSwitch.Fqdn + "|" + VSwitch.IPAdd
	mykey := conf.GetConfigItem("SWITCHID")
	log.Println("[PLANE][ANNOUNCELOCAL] Announcing  ", myannounce)

	myannounce_enc := crypt.FrameEncrypt([]byte(mykey), []byte(myannounce))

	tlv := tools.CreateTLV("A", myannounce_enc)

	DispatchTLV(tlv, mac)

}

func AnnounceAlien(alien_mac string, mac string) {

	mac = strings.ToUpper(mac)

	myannounce := strings.ToUpper(alien_mac) + "|" + VSwitch.Fqdn
	mykey := conf.GetConfigItem("SWITCHID")

	myannounce_enc := crypt.FrameEncrypt([]byte(mykey), []byte(myannounce))

	tlv := tools.CreateTLV("A", myannounce_enc)

	DispatchTLV(tlv, mac)

}

func SendQueryToMac(mac string) {

	mac = strings.ToUpper(mac)

	myannounce := VSwitch.HAddr
	mykey := conf.GetConfigItem("SWITCHID")

	myannounce_enc := crypt.FrameEncrypt([]byte(mykey), []byte(myannounce))

	tlv := tools.CreateTLV("Q", myannounce_enc)

	DispatchTLV(tlv, mac)

}
