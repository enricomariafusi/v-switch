package plane

import (
	"V-switch/conf"
	"log"
	"strconv"
)

var (
	TapToPlane chan []byte
	PlaneToTap chan []byte
	UdpToPlane chan []byte
)

func init() {

	queue_length, err := strconv.Atoi(conf.GetConfigItem("QUEUE"))

	if err != nil {
		queue_length = 256
	}

	TapToPlane = make(chan []byte, queue_length)
	UdpToPlane = make(chan []byte, queue_length)
	PlaneToTap = make(chan []byte, queue_length)

	log.Println("[PLANE][QUEUES] Queues created with lenght: ", queue_length)

}
