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
		log.Printf("[PLANE][TLV][INTERPRETER] Read %d bytes from UdpToPlane", len(my_tlv_enc))

		my_tlv := crypt.FrameDecrypt([]byte(VSwitch.SwID), my_tlv_enc)
		if my_tlv == nil {
			log.Printf("[PLANE][TLV][ERROR] Invalid KEY(%d): %s", len(VSwitch.SwID), VSwitch.SwID)
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
			if strings.Count(announce, "|") == 2 {
				fields := strings.Split(announce, "|")
				VSwitch.AddMac(fields[0], fields[1], fields[2])
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

func DispatchTLV(mytlv []byte, mac string) {

	mac = strings.ToUpper(mac)

	if VSwitch.macIsKnown(mac) {

		DispatchUDP(mytlv, VSwitch.SPlane[mac].EndPoint.String())

	}
}

func DispatchUDP(mytlv []byte, remote string) {

	var neterr error
	var RemoteAddr *net.UDPAddr

	var n int

	RemoteAddr, neterr = net.ResolveUDPAddr("udp", remote)
	if neterr != nil {
		log.Println("[PLANE][TLV][DispatchUDP] Remote address invalid :", neterr.Error())
		return
	}

	n, neterr = VSwitch.Server.WriteToUDP(mytlv, RemoteAddr) // we use the server IP and port as origin.
	if neterr != nil {
		log.Println("[PLANE][TLV][DispatchUDP] Error Writing to [", remote, "]:", neterr.Error())
		return
	} else {
		log.Printf("[PLANE][TLV][DispatchUDP] Written %d BYTES of %d to %s : %t", n, len(mytlv), remote, neterr == nil)
	}

}

func AnnounceLocal(mac string) {

	mac = strings.ToUpper(mac)

	myannounce := VSwitch.HAddr + "|" + VSwitch.Fqdn + "|" + VSwitch.IPAdd

	log.Printf("[PLANE][ANNOUNCELOCAL] Announcing [%s] ", myannounce)

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

	log.Printf("[PLANE][ANNOUNCEALIEN] Announcing [%s] ", myannounce)

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
