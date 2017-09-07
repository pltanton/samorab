package core

import (
	"log"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Listener struct {
	bot     *tgbotapi.BotAPI
	updates <-chan tgbotapi.Update
}

func InitListener(bot *tgbotapi.BotAPI) Listener {
	updatesChannel, err := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	if err != nil {
		panic("Can't subscribe to updates")
	}
	return Listener{
		bot:     bot,
		updates: updatesChannel,
	}
}

func (l Listener) Start() {
	for upd := range l.updates {
		log.Println(upd)
	}
}
