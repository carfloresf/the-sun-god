package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
)

const (
	configFileKey     = "configFile"
	defaultConfigFile = "config/config.yml"
	configFileUsage   = "config file path"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, configFileKey, defaultConfigFile, configFileUsage)
	flag.Parse()

	log.Println("configFile: ", configFile)

	execute(configFile)
}
