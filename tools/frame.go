package tools

import (
	"net"
	"strings"
)

func IsMacBcast(mac string) bool {

	m_hw, err := net.ParseMAC(mac)

	if err != nil {
		return false
	}

	return IsBroadcast(m_hw) || IsIPv4Multicast(m_hw)

}

func MACDestination(macFrame []byte) net.HardwareAddr {
	return net.HardwareAddr(macFrame[:6])
}

func MACSource(macFrame []byte) net.HardwareAddr {
	return net.HardwareAddr(macFrame[6:12])
}

func IsBroadcast(addr net.HardwareAddr) bool {
	return addr[0] == 0xff && addr[1] == 0xff && addr[2] == 0xff && addr[3] == 0xff && addr[4] == 0xff && addr[5] == 0xff
}

func IsIPv4Multicast(addr net.HardwareAddr) bool {
	return addr[0] == 0x01 && addr[1] == 0x00 && addr[2] == 0x5e
}

func CleanFrame(frame []byte) []byte {
	s := string(frame)
	s = strings.TrimRight(s, "\x00")
	return []byte(s)

}
