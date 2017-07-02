package tap

import (
	"log"

	"V-switch/plane"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func PcapReaderLoop() {

	var (
		snapshotLen int32 = int32(VDev.mtu * 2)
		err         error
		handle      *pcap.Handle
	)

	// Open device
	handle, err = pcap.OpenLive(VDev.devicename, snapshotLen, false, pcap.BlockForever)
	if err != nil {
		log.Printf("[TAP][PCAP] Cannot sniff the device  :%s(%d)", VDev.devicename, snapshotLen)

	} else {
		log.Printf("[TAP][PCAP] Start sniffing device <%s> with Snapshot Length %d bytes ", VDev.devicename, snapshotLen)
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
