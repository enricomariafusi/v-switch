package plane

import (
	"V-switch/conf"
	"V-switch/crypt"
	"V-switch/tools"
	"log"
	"net"
	"time"
)

func init() {

	if conf.ConfigItemExists("SEED") {
		seed_address := conf.GetConfigItem("SEED")
		log.Println("[PLANE][UDP][SEED]: Starting SEED to: ", seed_address)
		go SeedingTask(seed_address)
	} else {
		log.Println("[PLANE][UDP][SEED]: No SEED configured, not joining existing switch")
	}

}

func SeedingTask(remote string) {

	log.Println("[PLANE][UDP][SEED]: Creating conn with: ", remote)

	ServerAddr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		log.Println("[PLANE][UDP][SEED] Bad destination address ", remote, ":", err.Error())
		return
	}

	LocalAddr, err := net.ResolveUDPAddr("udp", tools.GetLocalIp()+":0")
	if err != nil {
		log.Println("[PLANE][UDP][SEED] Cannot find local port to bind ", remote, ":", err.Error())
		return
	}

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)

	if err != nil {
		log.Println("[PLANE][UDP][SEED] Error connecting with ", remote, ":", err.Error())
		return
	}
	log.Println("[PLANE][UDP][SEED] Success connecting with ", remote)
	mykey := conf.GetConfigItem("SWITCHID")

	for {

		// first, sends the announce

		var myfqdn string

		if conf.ConfigItemExists("PUBLIC") {
			myfqdn = conf.GetConfigItem("PUBLIC")
		} else {
			myfqdn = tools.GetFQDN() + ":" + conf.GetConfigItem("PORT")
			log.Println("[PLANE][UDP][SEED] dynamic hostid set to", myfqdn)
		}

		myannounce := VSwitch.HAddr + "|" + myfqdn

		myannounce_enc := crypt.FrameEncrypt([]byte(mykey), []byte(myannounce))

		tlv := tools.CreateTLV("A", myannounce_enc)

		_, err := Conn.Write(tlv)
		if err != nil {
			log.Println("[PLANE][UDP][SEED] Cannot announce to", remote)
		}

		// then sends query

		myannounce = VSwitch.HAddr

		myannounce_enc = crypt.FrameEncrypt([]byte(mykey), []byte(myannounce))

		tlv = tools.CreateTLV("Q", myannounce_enc)

		_, err = Conn.Write(tlv)
		if err != nil {
			log.Println("[PLANE][UDP][SEED] Cannot query to", remote)
		}

		time.Sleep(5 * time.Minute)

	}

}
