package main

import (
	"github.com/pltanton/samorab/configuration"
	"github.com/pltanton/samorab/core"
)

var listener core.Listener

func init() {
	configuration.InitFlag()
	configuration.InitConfig()
	bot := core.InitTelegramAPI()
	listener = core.InitListener(bot)
	listener.Start()
}

func main() {
}
