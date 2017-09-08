package main

import (
	"github.com/pltanton/samorab/configuration"
	"github.com/pltanton/samorab/core"
	"github.com/pltanton/samorab/storage"
)

var listener core.Listener
var synchronizer storage.DictionarySynchronizer

func init() {
	configuration.InitFlag()
	configuration.InitConfig()
	storage.InitStorage()
	bot := core.InitTelegramAPI()
	listener = core.InitListener(bot)
	synchronizer = storage.NewDictionarySynchronizer()
}

func main() {
	go synchronizer.Start()
	listener.Start()
}
