package core

import (
	"fmt"
	"math/rand"
	"time"
)

var messagesPool = []string{
	"Говорим как петербуржцы! Не %v, а %v!",
	"По-русски будет %[2]v, а ты всё %[1]v, да %[1]v!",
	"Правильно будет %[2]v, а ты все %[1]v, да %[1]v!",
	"Нет %v, а %v!",
	"Понаучатся в этих ваших интернетах! Не %v, а %v",
}

func formatMessage(from, to string) string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf(messagesPool[rand.Intn(len(messagesPool))], from, to)
}
