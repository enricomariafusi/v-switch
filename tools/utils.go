package tools

import (
	"bytes"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
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

func CreateTLV(typ string, payload string) []byte {

	tmp_length := strconv.Itoa(len(payload))

	var mybuffer bytes.Buffer

	mybuffer.WriteString(typ)
	mybuffer.WriteString(":")
	mybuffer.WriteString(tmp_length)
	mybuffer.WriteString(":")
	mybuffer.WriteString(payload)

	return mybuffer.Bytes()

}

func UnPackTLV(mytlv string) (typ string, ln int, payload string) {

	fields := strings.Split(mytlv, ":")

	if len(fields) != 3 {
		log.Println("[TOOLS][TLV] Invalid TLV : ", mytlv)
		return "N", 0, ""
	}

	plen, _ := strconv.Atoi(fields[1])

	if plen != len(fields[2]) {
		log.Println("[TOOLS][TLV] Malformed TLV ")
		return "N", 0, ""
	}

	return fields[0], plen, fields[2]

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
