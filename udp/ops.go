package netudp

import (
	"V-switch/conf"
	"V-switch/plane"
	"V-switch/tools"
	"log"
	"net"
	"os"
	"strconv"
)

func init() {

	myport := conf.GetConfigItem("PORT")

	if len(myport) < 2 {
		conf.SetConfigItem("PORT", "22000")
	}

	plane.VSwitch.Server = UdpCreateServer(conf.GetConfigItem("PORT"))
	go UDPReadMessage(plane.VSwitch.Server)

}

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		log.Println("[UDP] problem: ", err.Error())

	}
}

func UdpCreateServer(port string) *net.UDPConn {

	log.Println("[UDP][SERVER] Starting UDP listener")

	LocalAddr := tools.GetLocalIp() + ":" + port

	ServerAddr, err := net.ResolveUDPAddr("udp", LocalAddr)
	CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)

	// Setting the read buffer to avoid congestion

	ServerConn.SetReadBuffer(4194304)

	if err != nil {
		log.Println("[UDP][SERVER] Error listening at port ", port, ":", err.Error())
		os.Exit(1)
	}

	log.Println("[UDP][SERVER] Now listening at port ", port)
	return ServerConn

}

func UDPReadMessage(ServerConn *net.UDPConn) {

	log.Println("[UDP][SERVER] Reading thread started")

	tmp_MTU, _ := strconv.Atoi(conf.GetConfigItem("MTU")) // encryption of MTU max size

	defer ServerConn.Close()

	buf := make([]byte, 2*tmp_MTU) // enough for the payload , even if encrypted ang gob encoded
	log.Println("[UDP][SERVER] Read MTU set to ", tmp_MTU)

	for {

		n, addr, err := ServerConn.ReadFromUDP(buf)
		log.Println("[UDP][SERVER] Received ", n, "bytes from ", addr.String()) // just for debug

		if err != nil {
			log.Println("[UDP][SERVER] Error while reading: ", err.Error())
		} else {

			plane.UdpToPlane <- buf[:n]

		}

	}

}

func UDPEngineStart() {

	log.Println("[UDP] Engine init ")

}
