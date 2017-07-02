package tap

import (
	"log"

	"V-switch/plane"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"

	"time"
)

func PcapReaderLoop() {

	var (
		device      string = VDev.devicename
		snapshotLen int32  = int32(VDev.mtu)
		promiscuous bool   = false
		err         error
		timeout     time.Duration = 10 * time.Microsecond
		handle      *pcap.Handle
	)

	// Open device
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Printf("[TAP][PCAP] Cannot sniff the device  :%s(%d)", device, snapshotLen)

	} else {
		log.Printf("[TAP][PCAP] Start sniffing device <%s> with MTU %d ", device, snapshotLen)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		//fmt.Println("|----------------------------NEWPACKET-----------------------------------|")
		//fmt.Printf("Raw: % x \n", packet.Data())
		log.Printf("[TAP][PCAP][READ] Read %d Bytes long frame from TAP", len(packet.Data()))
		plane.TapToPlane <- packet.Data()
		log.Printf("[TAP][READ] Frame sent to Plane")
	}
}
