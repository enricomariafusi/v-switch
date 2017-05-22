package netudp

import (
	"V-switch/conf"
	"V-switch/tools"
	"log"
	"net"
	"strconv"
)

func init() {

	initport := "22000"

	if conf.ConfigItemExists("PORT") {

		initport := conf.GetConfigItem("PORT")
		log.Println("[UDP][CONF] Listening port: ", initport)
	} else {
		log.Println("[UDP][CONF] Port not configured, using 22000 ")
	}

	go UdpReceiveTLV(initport)

}

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		log.Println("[UDP] problem: ", err.Error())

	}
}

func UdpReceiveTLV(port string) net.Conn {

	log.Println("[UDP][SERVER] thread is on")

	tmp_MTU, _ := strconv.Atoi(conf.GetConfigItem("MTU")) // encryption of MTU max size

	tmp_MTU = (tmp_MTU + 16) * (4 / 3) // how much is the max payload by the base64 encoding overhead

	LocalAddr := tools.GetLocalIp() + ":" + port

	ServerAddr, err := net.ResolveUDPAddr("udp", LocalAddr)
	CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, tmp_MTU) // enough for the payload , even if encrypted and encoded base64

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		log.Println("Received ", string(buf[0:n]), " from ", addr.String()) // just for debug

		if err != nil {
			log.Println("[UDP][SERVER] Error while reading: ", err.Error())
		}

		// here the interpreter goes
		// here we check if we have the peer, and if we don't , we add it

	}
}

func UdpSendTLV(remote string, payload string) net.Conn {
	ServerAddr, err := net.ResolveUDPAddr("udp", remote)
	CheckError(err)

	Conn, err := net.Dial("udp", ServerAddr.String())
	CheckError(err)

	defer Conn.Close()

	for {

		//reads from the queue
		//encrypts
		//and writes like
		_, err := Conn.Write([]byte(payload))
		if err != nil {
			log.Println("[UDP][CLIENT]: Error sending payload: ", err.Error())
		}

	}
}

func UDPEngineStart() {

	log.Println("[UDP] Engine init ")

}
