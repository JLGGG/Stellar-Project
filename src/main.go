package main

import (
	"bytes"
	"fmt"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

// Token information is always kept private
const telegramBotToken string = ""

func makeAccount() (string, string, string) {
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}

	address := pair.Address()
	seed := pair.Seed()
	log.Printf("Secret key: %s", seed)
	log.Printf("Public key: %s", address)

	// Get 10,000 test XLM from friendbot.
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

	// Check account information
	request := horizonclient.AccountRequest{AccountID: address}
	account, err := horizonclient.DefaultTestNetClient.AccountDetail(request)
	if err != nil {
		log.Fatal(err)
	}

	var buffer bytes.Buffer

	log.Println("Balance for account:", address)
	// Used the bytes package for concatenating sentences. It's speed is O(n).
	buffer.WriteString(fmt.Sprintf("Account ID: https://horizon-testnet.stellar.org/accounts/%s\n", address))
	for _, balance := range account.Balances {
		log.Println(balance)
		buffer.WriteString(fmt.Sprintf("Account Balance: %s\n", balance.Balance))
	}

	log.Println(buffer.String())
	return address, seed, buffer.String()
}

func ParseBalanceStr(balanceStr string) string {
	// need to modify regular expression.
	regexp := regexp.MustCompile("{[0-9]+\\.[0-9]+")
	balanceStr = regexp.Split(balanceStr, 1)[0]

	return balanceStr
}

func main() {
	keysAndBalance := make([]string, 3)
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

	b.Handle("/make_account", func(m *tb.Message) {
		// keysAndBalance[0]: public key
		// keysAndBalance[1]: private key
		// keysAndBalance[2]: account's balance string
		keysAndBalance[0], keysAndBalance[1], keysAndBalance[2] = makeAccount()
		address := fmt.Sprintf("Public key: %s", keysAndBalance[0])
		seed := fmt.Sprintf("Secret key: %s", keysAndBalance[1])

		b.Send(m.Sender, address)
		b.Send(m.Sender, seed)
		b.Send(m.Sender, ParseBalanceStr(keysAndBalance[2]))
	})

	b.Start()
}
