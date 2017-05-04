package conf

import (
	ht "V-switch/tools"
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var VConfig map[string]string

func init() {

	VConfig = make(map[string]string)

	GConfigFile := filepath.Join(ht.Hpwd(), "etc", "vswitch.conf")

	readConfig(GConfigFile)

}

func StartConfig() {

	log.Printf("Reading config...\r\n")

}

func serializeConf(line string) {

	// create a splitter because "split" adds an empty line after the last \n
	splitter := func(c rune) bool {
		return (c == ' ' || c == '=') // trims space and understands equal
	}

	split := strings.FieldsFunc(line, splitter)

	if len(split) != 0 {

		VConfig[split[0]] = split[1]
		log.Printf("Config: %q -> %q\r\n", split[0], split[1])

	}

}

func readConfig(FileName string) {

	file, err := os.Open(FileName)
	if err != nil {
		log.Printf("[Config] can't open file %s", FileName)

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		serializeConf(line)
	}

	file.Close()

}
