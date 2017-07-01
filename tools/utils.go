package tools

import (
	"bytes"
	"encoding/base64"
	"log"
	"math/rand"
	"net"

	"strings"
	"time"
)

var letters = []rune("0123456789-+@abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//RandSeq returns a random string
func RandSeq(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//GetLocalIp returns back the IP of the interface hosting the default route
func GetLocalIp() string {
	// testing with  198.18.0.0/15 , see https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xhtml
	conn, err := net.Dial("udp", "198.18.0.30:80")
	if err != nil {
		log.Printf("[TOOLS][UTILS][OS] : cannot use UDP: %s", err.Error())
		return "127.0.0.1" // wanted to use 0.0.0.0 but golang cannot make use of it
	}
	conn.Close()
	torn := strings.Split(conn.LocalAddr().String(), ":")
	return torn[0]
}

func AddrResolve(fqdn string) (addr string) {

	addresses, err := net.LookupIP(fqdn)

	if err != nil {
		log.Printf("[TOOLS][UTILS][DNS] ERROR %s", err)
		// protocol = "tcp4"
		return "127.0.0.1"

	} else {

		addr := addresses[0].String()

		log.Printf("[TOOLS][UTILS][DNS] Resolution ok: %s -> %s", fqdn, addr)

		if strings.Count(addr, ":") > 2 {
			addr = "[" + addr + "]"
			// protocol = "tcp6"
			return addr
		}

		if strings.Contains(addr, ".") {
			// protocol = "tcp4"
			return addr
		}

	}

	return "127.0.0.1"

}

func CreateTLV(typ string, payload []byte) []byte {

	if len(payload) == 0 {
		return nil
	}

	var ret bytes.Buffer

	ret.WriteString(typ)
	ret.WriteString("|")
	ret.WriteString(base64.StdEncoding.EncodeToString(payload))

	return ret.Bytes()

}

func UnPackTLV(n_tlv []byte) (typ string, ln int, payload []byte) {

	if n_tlv[1] != '|' {
		log.Println("[TOOLS][TLV][UNPACK] WTF is this frame: ", string(n_tlv))
		return "Z", 0, nil
	}

	if strings.Contains("AQFD", string(n_tlv[0])) == false {
		log.Println("[TOOLS][TLV][UNPACK] WTF is this Type: ", string(n_tlv))
		return "Z", 0, nil
	}

	// n_tlv[0] contains the typ , tlv[2:] contains the payload

	tlvBin, err := base64.StdEncoding.DecodeString(string(n_tlv[2:]))

	if err != nil {
		log.Println("[TOOLS][TLV][UNPACK] Error with base64:", err.Error())
		return "Z", 0, nil
	}

	return string(n_tlv[0]), len(tlvBin), tlvBin

}

func GetFQDN() string {

	myIP := GetLocalIp()

	names, err := net.LookupAddr(myIP)
	if err != nil {
		log.Println("[TOOLS][UTILS][FQDN] Error getting my hostname:", err.Error())
		return "localhost"
	}

	return names[0]

}
