package main

import (
	"flag"
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

	execute(configFile)
}
