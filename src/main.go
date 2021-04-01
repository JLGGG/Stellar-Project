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

	//log.Println("Balance for account:", address)
	// Used the bytes package for concatenating sentences. It's speed is O(n).
	for _, balance := range account.Balances {
		//log.Println(balance)
		buffer.WriteString(fmt.Sprintf("%s\n", balance.Balance))
	}
	//TODO Id, Pw should be saved

	//log.Println(buffer.String())
	return address, seed, buffer.String()
}

func ParseBalanceStr(balanceStr string) string {
	// need to modify regular expression.
	regexp := regexp.MustCompile("[0-9]+\\.[0-9]+")
	balanceStr = regexp.FindAllString(balanceStr, 1)[0]

	return balanceStr
}

func main() {
	//TODO Add / command showing a list of commands.
	//	   Add remittance function.
	//     Add receive function.
	//     Add account view function.
	//     Add external API to get fiat money.
	//     Add anchor assets of XLM.
	//     Add function to get indicators such as USD index, 10 treasury, etc.

	keysAndBalance := make([]string, 3)
	var buffer bytes.Buffer
	b, err := tb.NewBot(tb.Settings{
		Token:  telegramBotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/", func(m *tb.Message) {
		b.Send(m.Sender, "List of Supported Commands:\n")
		b.Send(m.Sender, "/hello")
		b.Send(m.Sender, "/make_account")
	})

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello World!")
	})

	b.Handle("/make_account", func(m *tb.Message) {
		// keysAndBalance[0]: public key
		// keysAndBalance[1]: private key
		// keysAndBalance[2]: account's balance string
		keysAndBalance[0], keysAndBalance[1], keysAndBalance[2] = makeAccount()
		address := fmt.Sprintf("Public key(Id): %s", keysAndBalance[0])
		seed := fmt.Sprintf("Secret key(Pw): %s", keysAndBalance[1])

		balanceResult := ParseBalanceStr(keysAndBalance[2])
		buffer.WriteString(fmt.Sprintf("Account ID: https://horizon-testnet.stellar.org/accounts/%s\n", keysAndBalance[0]))
		buffer.WriteString(fmt.Sprintf("Current balance: %s\n", balanceResult))

		b.Send(m.Sender, address)
		b.Send(m.Sender, seed)
		b.Send(m.Sender, buffer.String())
	})

	b.Start()
}
