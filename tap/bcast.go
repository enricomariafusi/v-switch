package tap

import (
	"strings"
)

//Multicast address to be managed as broadcast into our switch
//ISO/IEC 9314-6 and ISO/IEC JTC1/SC25
var BCastMAC = []string{
	"01:00:0C:CC:CC:CC",
	"01:80:C2:00:00:02",
	"01:80:C2:00:00:1B",
	"01:80:C2:00:00:1C",
	"01:80:C2:00:00:1D",
	"01:80:C2:00:01:00",
	"01:00:5E:00:00:FB",
	"FF:FF:FF:FF:FF:FF"}

func IsMacBcast(mac string) bool {

	mac = strings.ToUpper(mac)

	for _, a := range BCastMAC {
		if strings.ToUpper(a) == mac {
			return true
		}
	}
	return false
}
