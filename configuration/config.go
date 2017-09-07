package configuration

import (
	"github.com/olebedev/config"
)

var cfg *config.Config

func GetCfg() *config.Config {
	return cfg
}

func InitConfig() {
	var err error
	cfg, err = config.ParseYamlFile(configPath)
	if err != nil {
		panic("Can't read configuration file" + configPath)
	}
}
