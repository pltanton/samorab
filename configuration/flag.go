package configuration

import (
	"flag"
)

var configPath string

func InitFlag() {
	flag.StringVar(&configPath, "config", "", "configuration file path")
	flag.Parse()

	if configPath == "" {
		panic("You should specify -config flag!")
	}
}
