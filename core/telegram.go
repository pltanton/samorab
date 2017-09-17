package core

import (
	"log"

	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/pltanton/samorab/configuration"
)

// InitTelegramAPI creates new BotAPI instance by settings from configulation
func InitTelegramAPI() *tgbotapi.BotAPI {
	key, error := configuration.GetCfg().String("telegram_api_key")
	if error != nil {
		panic("Can't read `telegram_api_key`")
	}

	bot, err := tgbotapi.NewBotAPI(key)
	if err != nil {
		panic("Can't initialize bot!")
	}
	log.Println("Bot successfully initialized")

	return bot
}
