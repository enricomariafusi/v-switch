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

		osocket := VSwitch.SPlane[mac].Socket

		_, err := osocket.Write([]byte(mytlv))
		log.Printf("[PLANE][TLV][DISPATCH] Sent %d BYTES to %s [%s]: %t", len(mytlv), mac, osocket.RemoteAddr().String(), err == nil)
		if err != nil {
			log.Printf("[PLANE][TLV][DISPATCH] cannot dispatch for MAC %s at [%s] : ", mac, osocket.RemoteAddr(), err.Error())
			VSwitch.RemoveMAC(mac)
		}

	} else {
		log.Println("[PLANE][TLV][DISPATCH] Unknown MAC: ", mac)

		return
	}

}

func CustomDispatch(mytlv []byte, remote string) {

	var neterr error
	var LocalAddr, RemoteAddr *net.UDPAddr
	var netsocket *net.UDPConn
	var n int

	LocalAddr, neterr = net.ResolveUDPAddr("udp", tools.GetLocalIp()+":0")
	if neterr != nil {
		log.Println("[PLANE][TLV][CustomDispatch] Error getting our IP :", neterr.Error())
		return
	}

	RemoteAddr, neterr = net.ResolveUDPAddr("udp", remote)
	if neterr != nil {
		log.Println("[PLANE][TLV][CustomDispatch] Remote address invalid :", neterr.Error())
		return
	}

	netsocket, neterr = net.DialUDP("udp", LocalAddr, RemoteAddr)
	if neterr != nil {
		log.Println("[PLANE][TLV][CustomDispatch] Error connecting with [", remote, "]:", neterr.Error())
		return
	}

	n, neterr = netsocket.Write(mytlv)
	if neterr != nil {
		log.Println("[PLANE][TLV][CustomDispatch] Error Writing to [", remote, "]:", neterr.Error())
		return
	} else {
		log.Printf("[PLANE][TLV][CustomDispatch] Written %d BYTES to %s : %t", n, remote, neterr == nil)
	}

	neterr = netsocket.Close()
	if neterr != nil {
		log.Println("[PLANE][TLV][CustomDispatch] Error closing socket : ", neterr.Error())
		return
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
