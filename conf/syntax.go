package conf

import (
	"log"
	"net"
	"os"
	"strconv"
)

func ConfCheck() {

	ConfigItems := []string{
		"MTU",
		"DEVICENAME",
		"PORT",
		"QUEUE",
		"SWITCHID",
	}

	OptionalItems := []string{
		"PUBLIC",
		"SEED",
	}

	// everything but 		"PUBLIC" and "SEED",

	for _, c_item := range ConfigItems {

		if !ConfigItemExists(c_item) {
			log.Printf("[CONF][SYNTAX] No %s in configuration. It is mandatory", c_item)
			os.Exit(1)
		}

	}

	// now some checks for syntax

	if _, err := strconv.Atoi(GetConfigItem("MTU")); err != nil {
		log.Println("[CONF][SYNTAX] Unacceptable value of MTU", GetConfigItem("MTU"))
		os.Exit(1)
	}

	if port, err := strconv.Atoi(GetConfigItem("PORT")); err != nil {
		log.Println("[CONF][SYNTAX] Unacceptable value of PORT", GetConfigItem("PORT"))
		os.Exit(1)
	} else {

		if (port > 65534) || (port < 10) {
			log.Println("[CONF][SYNTAX] Unacceptable value of PORT", GetConfigItem("PORT"))
			os.Exit(1)
		}
	}

	if queue, err := strconv.Atoi(GetConfigItem("QUEUE")); err != nil {
		log.Println("[CONF][SYNTAX] Unacceptable value of QUEUE", GetConfigItem("QUEUE"))
		os.Exit(1)
	} else {

		if (queue > 2048) || (queue < 16) {
			log.Println("[CONF][SYNTAX] Unacceptable value of QUEUE", GetConfigItem("QUEUE"))
			os.Exit(1)
		}
	}

	if len(GetConfigItem("SWITCHID")) != 32 {
		log.Println("[CONF][SYNTAX] Unacceptable value of SWITCHID", GetConfigItem("SWITCHID"))
		log.Println("[CONF][SYNTAX] It MUST be 32 char string with no spaces")
		os.Exit(1)
	}

	for _, c_item := range OptionalItems {

		saddr := GetConfigItem(c_item)

		if _, err := net.ResolveUDPAddr("udp", saddr); err != nil {
			log.Printf("[CONF][SYNTAX] %s is not a valid address for %s", saddr, c_item)
			os.Exit(1)

		}

	}

	if len(GetConfigItem("DEVICENAME")) > 14 {
		log.Println("[CONF][SYNTAX] Devicename too long", GetConfigItem("DEVICENAME"))
		os.Exit(1)
	}

	if len(GetConfigItem("DEVICENAME")) < 3 {
		log.Println("[CONF][SYNTAX] Devicename too short", GetConfigItem("DEVICENAME"))
		os.Exit(1)
	}

	log.Println("[CONF][SYNTAX] Conf syntax OK")

}
