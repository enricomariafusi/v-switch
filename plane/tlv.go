package plane

import (
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

	var my_tlv_enc []byte
	log.Println("[PLANE][TLV][INTERPRETER] Thread starts")

	for {

		my_tlv_enc = <-UdpToPlane

		my_tlv := crypt.FrameDecrypt([]byte(VSwitch.SwID), my_tlv_enc)
		if my_tlv == nil {
			continue
		}

		typ, ln, payload := tools.UnPackTLV(my_tlv)

		if ln == 0 {
			continue
		}

		switch typ {

		// it is a frame
		case "F":
			PlaneToTap <- payload
			// someone is announging itself
		case "A":
			announce := string(payload)
			fields := strings.Split(announce, "|")
			if len(fields) == 3 {
				radd, _ := net.ResolveUDPAddr("udp", fields[1])
				rip, _ := net.ResolveIPAddr("ip", fields[2])
				VSwitch.addPort(fields[0], *radd)
				UDPCreateConn(fields[0], *radd)
				VSwitch.addRemoteIp(fields[0], *rip)
			}
		case "Q":
			sourcemac := string(payload)
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

	log.Println("[PLANE][ANNOUNCELOCAL] Announcing  ", myannounce)

	tlv := tools.CreateTLV("A", []byte(myannounce))

	tlv_enc := crypt.FrameEncrypt([]byte(VSwitch.SwID), tlv)

	DispatchTLV(tlv_enc, mac)

}

// Announces  port which is not ours
func AnnounceAlien(alien_mac string, mac string) {

	mac = strings.ToUpper(mac)
	alien_mac = strings.ToUpper(alien_mac)

	tmp_endpoint := VSwitch.SPlane[alien_mac].EndPoint
	tmp_ethIP := VSwitch.SPlane[alien_mac].EthIP

	myannounce := alien_mac + "|" + tmp_endpoint.String() + "|" + tmp_ethIP.String()

	tlv := tools.CreateTLV("A", []byte(myannounce))

	tlv_enc := crypt.FrameEncrypt([]byte(VSwitch.SwID), tlv)

	DispatchTLV(tlv_enc, mac)

}

func SendQueryToMac(mac string) {

	mac = strings.ToUpper(mac)

	myannounce := VSwitch.HAddr

	tlv := tools.CreateTLV("Q", []byte(myannounce))

	tlv_enc := crypt.FrameEncrypt([]byte(VSwitch.SwID), tlv)

	DispatchTLV(tlv_enc, mac)

}
