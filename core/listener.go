package core

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/un000/mystem-wrapper"
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/pltanton/samorab/configuration"
	"github.com/pltanton/samorab/storage"
)

var mystem *mystem_wrapper.MyStem

// Listener watches of telegram updates and process it
type Listener struct {
	bot     *tgbotapi.BotAPI
	updates <-chan tgbotapi.Update
}

// InitListener initializes new Listener with given telegram BotAPI
func InitListener(bot *tgbotapi.BotAPI) Listener {
	mystemBin, err := configuration.GetCfg().String("mystem_bin")
	if err != nil {
		log.Fatalf("Can't find mystem in conf file")
	}
	mystem = mystem_wrapper.New(mystemBin, []string{})
	updatesChannel, err := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	if err != nil {
		panic("Can't subscribe to updates")
	}
	return Listener{
		bot:     bot,
		updates: updatesChannel,
	}
}

// Start starts monitoring of telegram updates
func (l Listener) Start() {
	for upd := range l.updates {
		l.processUpdate(upd)
	}
}

func (l Listener) processUpdate(upd tgbotapi.Update) {
	rand.Seed(time.Now().Unix())

	if upd.Message != nil {
		switch {
		case upd.Message.IsCommand():
			go l.processCommand(upd.Message)
		case rand.Intn(100) <= storage.GetChance(int(upd.Message.Chat.ID)):
			go l.processMessage(upd.Message)
		}
	}
}

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
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyToMessageID = message.MessageID
	l.bot.Send(msg)
}

func (l Listener) processCommand(message *tgbotapi.Message) {
	log.Printf("Recieved a command `%v` from chat #%v\n", message.Command(), message.Chat.ID)
	switch command := message.Command(); command {
	case "verojatnost":
		l.commandChance(message)
	case "perevedi":
		l.commandTranslate(message)
	}
}

func (l Listener) commandChance(message *tgbotapi.Message) {
	if message.CommandArguments() == "" {
		currentChance := storage.GetChance(int(message.Chat.ID))
		l.replyToMessage(message, fmt.Sprintf("Current chance to reply: %v%%", currentChance))
	} else {
		currentChance := storage.GetChance(int(message.Chat.ID))
		argument, err := strconv.ParseInt(message.CommandArguments(), 10, 32)
		if err != nil || argument <= 0 || argument > 100 {
			l.replyToMessage(message, "Ты своей головой сам подумай то! Вероятность должна быть в промежутке от 1 до 100!")
			return
		}
		storage.SetChance(int(message.Chat.ID), int(argument))
		log.Printf("Set chance %v for chat #%v", argument, message.Chat.ID)

		l.replyToMessage(message, fmt.Sprintf("Раньше я мог ответить с вероятностью %v%%, атеперь с %v%%!", currentChance, argument))
	}
}

func (l Listener) commandTranslate(message *tgbotapi.Message) {
	alternatives := storage.FindAlternatives(message.CommandArguments())
	if alternatives == nil {
		l.replyToMessage(message, "Кажется, я ничего не знаю об этом слове.")
		return
	}
	words := append(alternatives.Alternative, alternatives.Synonim...)
	wordsBold := make([]string, len(words))
	for i, word := range words {
		wordsBold[i] = fmt.Sprintf("__%v__", word)
	}
	l.replyToMessage(
		message,
		fmt.Sprintf("Я бы заменил это слово на: %v", strings.Join(wordsBold, " или ")),
	)
}
