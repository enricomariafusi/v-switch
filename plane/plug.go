package plane

import (
	"V-switch/conf"
	"V-switch/crypt"
	"V-switch/tools"
	"log"
	"net"
	"strconv"
	"time"
)

func init() {

	if conf.ConfigItemExists("SEED") {
		seed_address := conf.GetConfigItem("SEED")
		log.Println("[PLANE][PLUG]: Starting SEED to: ", seed_address)
		go SeedingTask(seed_address)
	} else {
		log.Println("[PLANE][PLUG]: No SEED configured, not joining existing switch")
	}

}

func SeedingTask(remote string) {

	log.Println("[PLANE][PLUG]: Creating conn with: ", remote)

	ServerAddr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		log.Println("[PLANE][PLUG] Bad destination address ", remote, ":", err.Error())
		return
	}

	LocalAddr, err := net.ResolveUDPAddr("udp", tools.GetLocalIp()+":0")
	if err != nil {
		log.Println("[PLANE][PLUG] Cannot find local port to bind ", remote, ":", err.Error())
		return
	}

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)

	if err != nil {
		log.Println("[PLANE][PLUG] Error connecting with ", remote, ":", err.Error())
		return
	}
	log.Println("[PLANE][PLUG] Success connecting with ", remote)
	mykey := conf.GetConfigItem("SWITCHID")

	for {

		// first, sends the announce

		var myfqdn string

		if conf.ConfigItemExists("PUBLIC") {
			myfqdn = conf.GetConfigItem("PUBLIC")
		} else {
			myfqdn = tools.GetFQDN() + ":" + conf.GetConfigItem("PORT")
			log.Println("[PLANE][PLUG] dynamic hostid set to", myfqdn)
		}

		myannounce := VSwitch.HAddr + "|" + myfqdn

		myannounce_enc := crypt.FrameEncrypt([]byte(mykey), []byte(myannounce))

		tlv := tools.CreateTLV("A", myannounce_enc)

		_, err := Conn.Write(tlv)
		if err != nil {
			log.Printf("[PLANE][PLUG] Cannot announce to %s : %s", remote, err.Error())
		}

		// then sends query

		myannounce = VSwitch.HAddr

		myannounce_enc = crypt.FrameEncrypt([]byte(mykey), []byte(myannounce))

		tlv = tools.CreateTLV("Q", myannounce_enc)

		_, err = Conn.Write(tlv)
		if err != nil {
			log.Printf("[PLANE][PLUG] Cannot query to %s: %s", remote, err.Error())
		}

		cycle, _ := strconv.Atoi(conf.GetConfigItem("TTL"))

		time.Sleep(time.Duration(cycle) * time.Second)

	}

}
