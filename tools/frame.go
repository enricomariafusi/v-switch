package tools

import (
	"bytes"
	"net"
	"strings"
)

var BCastMAC = []string{
	"01:00:0C:CC:CC:CC",
	"01:00:0C:CC:CC:CD",
	"01:80:C2:00:00:00",
	"01:80:C2:00:00:02",
	"01:80:C2:00:00:1B",
	"01:80:C2:00:00:1C",
	"01:80:C2:00:00:1D",
	"01:80:C2:00:01:00",
	"FF:FF:FF:FF:FF:FF"}

func IsMacBcast(mac string) bool {

	m_hw, err := net.ParseMAC(mac)

	if err != nil {
		return false
	}

	return IsBroadcast(m_hw) || IsIPMulticast(m_hw)

}

func MACDestination(macFrame []byte) net.HardwareAddr {
	return net.HardwareAddr(macFrame[:6])
}

func MACSource(macFrame []byte) net.HardwareAddr {
	return net.HardwareAddr(macFrame[6:12])
}

func IsBroadcast(addr net.HardwareAddr) bool {
	return IsTierBcast(addr.String())
}

func IsIPv4Multicast(addr net.HardwareAddr) bool {
	return addr[0] == 0x01 && addr[1] == 0x00 && addr[2] == 0x5e
}

func IsIPv6Multicast(addr net.HardwareAddr) bool {
	return addr[0] == 0x33 && addr[1] == 0x33
}

func IsIPMulticast(addr net.HardwareAddr) bool {
	return IsIPv4Multicast(addr) || IsIPv6Multicast(addr)
}

func CleanFrame(frame []byte) []byte {

	splitter := func(c rune) bool {
		return (c == 0) //
	}

	return bytes.TrimRightFunc(frame, splitter)

}

func IsTierBcast(mac string) bool {

	mac = strings.ToUpper(mac) //we don't need to repeat for each element

	for _, a := range BCastMAC {
		if strings.ToUpper(a) == mac {
			return true
		}
	}
	return false
}
