package plane

import (
	"V-switch/conf"
	"log"
	"strconv"
)

type NetMessage struct {
	ETlv []byte // the payload (a whole encrypted TLV)
	Addr string // IP:PORT of the sender

}

var (
	TapToPlane chan []byte
	PlaneToTap chan []byte
	UdpToPlane chan NetMessage
)

func init() {

	queue_length, err := strconv.Atoi(conf.GetConfigItem("QUEUE"))

	if err != nil {
		queue_length = 256
	}

	TapToPlane = make(chan []byte, queue_length)
	UdpToPlane = make(chan NetMessage, queue_length)
	PlaneToTap = make(chan []byte, queue_length)

	log.Println("[PLANE][QUEUES] Queues created with lenght: ", queue_length)

}
