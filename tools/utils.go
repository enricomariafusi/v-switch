package tools

import (
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

	conn, err := net.Dial("udp", "255.255.255.255:80")
	if err != nil {
		log.Printf("[TOOLS] SYSADMIIIIIN : cannot use UDP")
		return "127.0.0.1" // wanted to use 0.0.0.0 but golang didn't get this
	}
	conn.Close()
	torn := strings.Split(conn.LocalAddr().String(), ":")
	return torn[0]
}
