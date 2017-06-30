package netudp

import (
	"V-switch/conf"
	"V-switch/plane"
	"V-switch/tools"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

var NetM plane.NetMessage

func init() {

	go UDPReadMessage()

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

	log.Println("[UDP][SERVER] Now listening at:  ", ServerConn.LocalAddr().String())
	return ServerConn

}

func UDPReadMessage() {

	defer func() {
		if e := recover(); e != nil {
			log.Println("[UDP][SERVER] Network listener issue, trying to save it")
			err, ok := e.(error)
			if !ok {
				err = fmt.Errorf("[EXC]: %v", e)
			}
			log.Printf("[UDP][SERVER] Error: <%s>", err)

		}
	}()

	plane.VSwitch.Server = UdpCreateServer(conf.GetConfigItem("PORT"))
	ServerConn := plane.VSwitch.Server

	log.Println("[UDP][SERVER] Reading thread started")

	readbuffer, _ := strconv.Atoi(conf.GetConfigItem("MTU")) // at least the  MTU max size

	buf := make([]byte, 3*readbuffer) // enough for the payload , even if encrypted ang gob encoded
	log.Println("[UDP][SERVER] Read MTU set to ", 3*readbuffer)

	for {

		n, addr, err := ServerConn.ReadFromUDP(buf)
		log.Println("[UDP][SERVER] Received ", n, "bytes from ", addr.String()) // just for debug

		if err != nil {
			log.Println("[UDP][SERVER] Error while reading: ", err.Error())
		} else {

			go func(msg []byte, ind *net.UDPAddr) {
				NetM.ETlv = msg
				NetM.Addr = ind.String()
				plane.UdpToPlane <- NetM

			}(buf[:n], addr)

		}

	}

}

func UDPEngineStart() {

	log.Println("[UDP] Engine init ")

}
