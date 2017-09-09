package core

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
	"regexp"
	"strings"

	"github.com/un000/mystem-wrapper"
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/pltanton/samorab/storage"
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
		l.processUpdate(upd)
	}
}

func (l Listener) processUpdate(upd tgbotapi.Update) {
	rand.Seed(time.Now().Unix())

	if upd.Message != nil {
		if upd.Message.IsCommand() {
			go l.processCommand(upd.Message)
		}
		chance := storage.GetChance(int(upd.Message.Chat.ID))
		if rand.Intn(100) <= chance {
			go l.processMessage(upd.Message)
		}
	}
}

var mystem = mystem_wrapper.New("/usr/bin/mystem", []string{})


func (l Listener) processMessage(message *tgbotapi.Message) {
	rand.Seed(time.Now().Unix())

	reg, _ := regexp.Compile("[^а-яА-ЯёЁ]")
	rawWords := strings.Fields(reg.ReplaceAllString(message.Text, " "))

	words, err := mystem.Transform(rawWords)
	if err != nil {
		log.Printf("Can't mystem: %v\n", message.Text)
		return
	}

	records := make([]*storage.DictionaryRecord, 0)
	for _, word := range words {
		if word == "" {
			continue
		}

		record := storage.FindAlternatives(word)
		if record == nil {
			continue
		}
		records = append(records, record)
	}
	if len(records) == 0 {
		return
	}

	record := records[rand.Intn(len(records))]

	words = append(record.Alternative, record.Synonim...)
	if len(words) == 0 {
		return
	}
	replacement := words[rand.Intn(len(words))]

	l.replyToMessage(message, formatMessage(record.Original, replacement))
}

func (l Listener) replyToMessage(message *tgbotapi.Message, answer string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, answer)
	msg.ReplyToMessageID = message.MessageID
	l.bot.Send(msg)
}

func (l Listener) processCommand(message *tgbotapi.Message) {
	log.Printf("Recieved a command `%v` from chat #%v\n", message.Command(), message.Chat.ID)
	switch command := message.Command(); command {
	case "verojatnost":
		if message.CommandArguments() == "" {
			currentChance := storage.GetChance(int(message.Chat.ID))
			l.replyToMessage(message, fmt.Sprintf("Current chance to reply: %v%%", currentChance))
		} else {
			currentChance := storage.GetChance(int(message.Chat.ID))
			argument, err := strconv.ParseInt(message.CommandArguments(), 10, 32)
			if err != nil || argument <= 0 || argument > 100 {
				l.replyToMessage(message, "I can set chance only in range of 1 to 100, stop bullshitting me!")
				return
			}
			storage.SetChance(int(message.Chat.ID), int(argument))
			log.Printf("Set chance %v for chat #%v", argument, message.Chat.ID)

			l.replyToMessage(message, fmt.Sprintf("Reply chance changed from %v%% to %v%%", currentChance, argument))
		}
	}
}
