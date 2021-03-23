package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

// Token information is always kept private
const telegramBotToken string = ""

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  telegramBotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello World!")
	})

	b.Start()
}
