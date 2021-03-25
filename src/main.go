package main

import (
	"fmt"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"net/http"
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

	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + address)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))

	request := horizonclient.AccountRequest{AccountID: address}
	account, err := horizonclient.DefaultTestNetClient.AccountDetail(request)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Balance for account:", address)

	for _, balance := range account.Balances {
		log.Println(balance)
	}
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
