package conf

import (
	ht "V-switch/tools"
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//VConfig will export the configuration to any package needing it
var VConfig map[string]string

func init() {

	VConfig = make(map[string]string)

	GConfigFile := filepath.Join(ht.Hpwd(), "etc", "vswitch.conf")

	readConfig(GConfigFile)

	ConfCheck()

}

//StartConfig is used in main.init to trigger the init in the package conf
func StartConfig() {

	log.Println("[CONF] Config Engine Init")

}

func serializeConf(line string) {

	// create a splitter because "split" adds an empty line after the last \n
	splitter := func(c rune) bool {
		return (c == ' ' || c == '=') // trims space and understands equal
	}

	split := strings.FieldsFunc(line, splitter)

	if len(split) != 0 {

		VConfig[split[0]] = split[1]
		log.Printf("[CONF]: %q -> %q\r\n", split[0], split[1])

	}

}

func readConfig(FileName string) {

	file, err := os.Open(FileName)
	if err != nil {
		log.Printf("[CONF] can't open file %s", FileName)

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		serializeConf(line)
	}

	file.Close()

}

func ConfigItemExists(item string) bool {

	_, exists := VConfig[item]

	return exists

}

func GetConfigItem(item string) string {

	if ConfigItemExists(item) {
		return VConfig[item]
	} else {
		return ""
	}

}

func SetConfigItem(item string, value string) {

	VConfig[item] = value

}
