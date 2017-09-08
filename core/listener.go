package core

import (
	//"log"
	"math/rand"
	"time"

	"github.com/kljensen/snowball"
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/pltanton/samorab/storage"
	"github.com/pltanton/samorab/utils"
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
		l.processMessage(upd)
	}
}

func (l Listener) processMessage(upd tgbotapi.Update) {
	rand.Seed(time.Now().Unix())

	if upd.Message == nil {
		return
	}
	records := make([]*storage.DictionaryRecord, 0)
	for _, word := range utils.SPLIT_REGEX.Split(upd.Message.Text, -1) {
		if word == "" {
			continue
		}
		stemmedWord, _ := snowball.Stem(word, "russian", true)
		record := storage.FindAlternatives(stemmedWord)
		if record == nil {
			continue
		}
		records = append(records, record)
	}
	if len(records) == 0 {
		return
	}

	record := records[rand.Intn(len(records))]

	words := append(record.Alternative, record.Synonim...)
	if len(words) == 0 {
		return
	}
	replacement := words[rand.Intn(len(words))]

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, formatMessage(record.Original, replacement))
	msg.ReplyToMessageID = upd.Message.MessageID
	l.bot.Send(msg)
}
