package tools

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"

	"strings"
	"time"
)

type TlvJ struct {
	T string
	P []byte
}

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

	var MyTlv TlvJ

	MyTlv.T = typ
	MyTlv.P = payload

	b, err := json.Marshal(MyTlv)
	if err != nil {
		log.Println("[TOOLS][TLV][CREATE] TLV creation failed: ", err.Error())
		return nil
	}

	log.Println("[TOOLS][TLV][CREATE] Frame created: ", string(b))

	return b

}

func UnPackTLV(n_tlv []byte) (typ string, ln int, payload []byte) {

	var MyTlv TlvJ

	err := json.Unmarshal(n_tlv, &MyTlv)
	if err != nil {
		log.Println("[TOOLS][TLV][UNPACK] TLV Unpack failed: ", err.Error())
		log.Println("[TOOLS][TLV][UNPACK] TLV Offending: ", string(n_tlv))
		return "Z", 0, nil
	}

	return MyTlv.T, len(MyTlv.P), MyTlv.P

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
