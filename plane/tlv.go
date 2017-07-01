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

	log.Println("[PLANE][TLV][INTERPRETER] Thread starts")

	for my_tlv_enc := range UdpToPlane {

		go interpreter(my_tlv_enc)

	}

}

func interpreter(mytlvenc NetMessage) {

	log.Printf("[PLANE][TLV][INTERPRETER] Read %d bytes from UdpToPlane", len(mytlvenc.ETlv))

	my_tlv := crypt.FrameDecrypt([]byte(VSwitch.SwID), mytlvenc.ETlv)
	if my_tlv == nil {
		log.Printf("[PLANE][TLV][ERROR] Invalid KEY(%d): %s", len(VSwitch.SwID), VSwitch.SwID)
		return
	} else {
		log.Printf("[PLANE][TLV][INTERPRETER] Decrypted JSON %d BYTES long", len(my_tlv))
	}

	typ, ln, payload := tools.UnPackTLV(my_tlv)

	if ln == 0 {
		log.Printf("[PLANE][TLV][ERROR] Payload was empty, nothing to do")
		return
	}

	log.Println("[PLANE][TLV][INTERPRETER] Received valid payload, type [", typ, "]")

	switch typ {

	// it is a frame
	case "F":
		PlaneToTap <- payload
		// someone is announging itself
	case "A":
		announce := string(payload)
		if strings.Count(announce, "|") == 1 {
			fields := strings.Split(announce, "|")
			VSwitch.AddMac(fields[0], mytlvenc.Addr, fields[1])
		}

	case "D":
		announce := string(payload)
		if strings.Count(announce, "|") == 2 {
			fields := strings.Split(announce, "|")
			VSwitch.AddMac(fields[0], fields[1], fields[2])
		}

	case "Q":
		sourcemac := string(payload)

		answerquery(sourcemac)

	default:
		log.Println("[PLANE][TLV][INTERPRETER] Unknown type, discarded: [ ", typ, " ]")

	}

}

func answerquery(mac string) {

	if _, err := net.ParseMAC(mac); err != nil {
		log.Printf("[PLANE][TLV][QUERY] Invalid mac %s : %s", mac, err.Error())
		return
	}

	AnnounceLocal(mac)

	if len(VSwitch.SPlane) > 0 {
		for alienmac := range VSwitch.SPlane {

			AnnounceAlien(alienmac, mac)
		}
		log.Println("[PLANE][TLV][QUERY] Query answered with success")
	} else {
		log.Println("[PLANE][TLV][QUERY] PLANE EMPTY, NO ANSWER TO QUERY")

	}

}

func DispatchTLV(mytlv []byte, mac string) {

	mac = strings.ToUpper(mac)

	if VSwitch.macIsKnown(mac) {

		DispatchUDP(mytlv, VSwitch.SPlane[mac].EndPoint.String())

	} else {
		log.Println("[PLANE][TLV][DISPATCH] Unknown MAC : [ ", mac, " ]")
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

	var strs []string
	myannounce := strings.Join(append(strs, VSwitch.HAddr, VSwitch.IPAdd), "|")

	log.Printf("[PLANE][ANNOUNCELOCAL] Announcing [%s] ", myannounce)

	tlv := tools.CreateTLV("A", []byte(myannounce))

	tlv_enc := crypt.FrameEncrypt([]byte(VSwitch.SwID), tlv)

	DispatchTLV(tlv_enc, mac)

}

// Announces  port which is not ours
func AnnounceAlien(alien_mac string, mac string) {

	mac = strings.ToUpper(mac)
	alien_mac = strings.ToUpper(alien_mac)

	var strs []string
	myannounce := strings.Join(append(strs, alien_mac, VSwitch.SPlane[alien_mac].EndPoint.String(), VSwitch.SPlane[alien_mac].EthIP.String()), "|")

	log.Printf("[PLANE][ANNOUNCEALIEN] Announcing [%s] ", myannounce)

	tlv := tools.CreateTLV("D", []byte(myannounce))

	tlv_enc := crypt.FrameEncrypt([]byte(VSwitch.SwID), tlv)

	DispatchTLV(tlv_enc, mac)

}

func SendQueryToMac(mac string) {

	mac = strings.ToUpper(mac)

	myannounce := VSwitch.HAddr

	tlv := tools.CreateTLV("Q", []byte(VSwitch.HAddr))

	tlv_enc := crypt.FrameEncrypt([]byte(VSwitch.SwID), tlv)

	log.Printf("[PLANE][ANNOUNCEALIEN] Querying %s with our mac %s ", mac, myannounce)

	DispatchTLV(tlv_enc, mac)

}
