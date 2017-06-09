package tools

import (
	"bytes"
	"encoding/gob"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

var letters = []rune("0123456789-+@abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Tlv struct {
	t       string // type of TLV
	payload []byte
}

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

	conn, err := net.Dial("udp", "255.255.255.255:80")
	if err != nil {
		log.Printf("[TOOLS][OS] : cannot use UDP")
		return "127.0.0.1" // wanted to use 0.0.0.0 but golang didn't get this
	}
	conn.Close()
	torn := strings.Split(conn.LocalAddr().String(), ":")
	return torn[0]
}

func AddrResolve(fqdn string) (addr string) {

	addresses, err := net.LookupIP(fqdn)

	if err != nil {
		log.Printf("[TOOLS][DNS] ERROR %s", err)
		// protocol = "tcp4"
		return "127.0.0.1"

	} else {

		addr := addresses[0].String()

		log.Printf("[TOOLS][DNS] Resolution ok: %s -> %s", fqdn, addr)

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

	var mybuffer bytes.Buffer

	encoder := gob.NewEncoder(&mybuffer)

	err := encoder.Encode(Tlv{typ, payload})
	if err != nil {
		log.Println("[TOOLS][TLV] Problem encoding")
	}

	return mybuffer.Bytes()

}

func UnPackTLV(n_tlv []byte) (typ string, ln int, payload []byte) {

	var mytlv Tlv
	var mybuffer bytes.Buffer

	decoder := gob.NewDecoder(&mybuffer)

	mybuffer.Write(n_tlv)
	err := decoder.Decode(&mytlv)
	if err != nil {
		log.Println("[TOOLS][TLV] Error recoding TLV")
		return "", 0, nil

	}

	return mytlv.t, len(mytlv.payload), mytlv.payload

}

func GetFQDN() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "localhost"
	}

	host := AddrResolve(hostname)
	if host == "127.0.0.1" {
		return "localhost"
	}

	return hostname
}
