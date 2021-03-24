package main

import (
	"fmt"
	"github.com/stellar/go/keypair"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

// Token information is always kept private
const telegramBotToken string = ""

func makeAccount() (string, string) {
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}

	address := pair.Address()
	seed := pair.Seed()
	log.Printf("Secret key: %s", seed)
	log.Printf("Public key: %s", address)
	return address, seed
}

func main() {
	keys := make([]string, 2)
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

	// need to add key-value DB for saving keys at this point.
	b.Handle("/make_account", func(m *tb.Message) {
		keys[0], keys[1] = makeAccount()
		address := fmt.Sprintf("Public key: %s", keys[0])
		seed := fmt.Sprintf("Secret key: %s", keys[1])
		b.Send(m.Sender, address)
		b.Send(m.Sender, seed)
	})

	b.Start()
}
