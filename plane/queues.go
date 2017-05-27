package plane

import (
	"V-switch/conf"
	"log"
	"strconv"

	"github.com/songgao/packets/ethernet"
)

var (
	TapToPlane chan ethernet.Frame
	PlaneToTap chan ethernet.Frame
	UdpToPlane chan string
	PlaneToUdp chan string
)

func init() {

	queue_length, err := strconv.Atoi(conf.GetConfigItem("QUEUE"))

	if err != nil {
		queue_length = 256
	}

	TapToPlane = make(chan ethernet.Frame, queue_length)
	UdpToPlane = make(chan string, queue_length)
	PlaneToTap = make(chan ethernet.Frame, queue_length)
	PlaneToUdp = make(chan string, queue_length)
	log.Println("[PLANE][QUEUES] Queues created with lenght: ", queue_length)

}
