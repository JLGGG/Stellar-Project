package main

import (
	_ "bufio"
	"bytes"
	"fmt"
	"github.com/JLGGG/Stellar-Project/src/stellar"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	_ "strconv"
	"time"
)

// Token information is always kept private
const telegramBotToken string = ""

// Save file name
const fileNameAboutEntity string = "src/info.txt"

func main() {
	// TODO:
	//	   Add remittance function.
	//	   Add account favorites
	//     Add receive function.
	//     Add external API to get fiat money.
	//     Add anchor assets of XLM.
	//     Add function to get indicators such as USD index, 10 treasury, etc.

	keysAndBalance := make([]string, 3)
	var buffer bytes.Buffer
	b, err := tb.NewBot(tb.Settings{
		Token:  telegramBotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	stellar.CheckError(err)

	b.Handle("/", func(m *tb.Message) {
		b.Send(m.Sender, "List of Supported Commands:\n")
		b.Send(m.Sender, "/hello")
		b.Send(m.Sender, "/make_account")
		b.Send(m.Sender, "/show_account")
		b.Send(m.Sender, "/send_payment")
	})

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello World!")
	})

	//TODO What to do when /make_account called twice?
	b.Handle("/make_account", func(m *tb.Message) {
		// keysAndBalance[0]: public key
		// keysAndBalance[1]: private key
		// keysAndBalance[2]: account's balance string
		keysAndBalance[0], keysAndBalance[1], keysAndBalance[2] = stellar.MakeAccount()
		address := fmt.Sprintf("Public key(Id): %s\n", keysAndBalance[0])
		seed := fmt.Sprintf("Secret key(Pw): %s\n", keysAndBalance[1])

		balanceResult := stellar.ParseBalanceStr(keysAndBalance[2])
		buffer.WriteString(fmt.Sprintf("Account ID: https://horizon-testnet.stellar.org/accounts/%s\n", keysAndBalance[0]))
		buffer.WriteString(fmt.Sprintf("Current balance: %s\n", balanceResult))

		//var temp []byte
		//temp = append(temp, []byte(address)...)
		//temp = append(temp, []byte(seed)...)
		//temp = append(temp, []byte(buffer.String())...)
		//writeFile(fileNameAboutEntity, temp)
		stellar.WriteTxToDB(keysAndBalance[0], keysAndBalance[1])

		b.Send(m.Sender, address)
		b.Send(m.Sender, seed)
		b.Send(m.Sender, buffer.String())
	})

	b.Handle("/show_account", func(m *tb.Message) {
		//b.Send(m.Sender, string(readFile(fileNameAboutEntity)))
		sID := make([]string, 100)
		sPW := make([]string, 100)
		count := stellar.ReadTxFromDB(sID, sPW)

		b.Send(m.Sender, "View all current account information")
		for i := 0; i < count; i++ {
			strID := fmt.Sprintf("Public key(Id): %s\n", sID[i])
			strPW := fmt.Sprintf("Secret key(Pw): %s\n", sPW[i])
			note := fmt.Sprintf("------------------ Account number: %d------------------", i+1)
			b.Send(m.Sender, note)
			b.Send(m.Sender, strID)
			b.Send(m.Sender, strPW)
		}
	})

	b.Handle("/send_payment", func(m *tb.Message) {
		// Enter the address to send
		// Check frequently used accounts
		b.Send(m.Sender, fmt.Sprintln("This is a remittance command."))
		b.Send(m.Sender, fmt.Sprint("Would you like to browse your frequently used accounts? (y/n): "))
		b.Handle(tb.OnText, func(m *tb.Message) {
			if m.Text == "y" {
				sID := make([]string, 100)
				count := stellar.ReadFavoriteAccountFromDB(sID)

				b.Send(m.Sender, "The list of favorite accounts:")
				for i := 0; i < count; i++ {
					b.Send(m.Sender, fmt.Sprintf("------------------ Account number: %d------------------", i+1))
					b.Send(m.Sender, fmt.Sprintf("Public key(Id): %s\n", sID[i]))
				}
			} else {
				var src, dst string
				sID := make([]string, 100)
				sPW := make([]string, 100)
				b.Send(m.Sender, fmt.Sprint("Enter the address to send(Press \"list\" to see the list of accounts):"))
				b.Handle(tb.OnText, func(m *tb.Message) {
					if len(m.Text) != 56 && m.Text != "list" {
						b.Send(m.Sender, "Account ID is 56 characters. Please re-enter.")
					} else if m.Text == "list" {
						b.Send(m.Sender, "Please select the account you want to use:")

						count := stellar.ReadTxFromDB(sID, sPW)
						for i := 0; i < count; i++ {
							b.Send(m.Sender, fmt.Sprintf("------------------ Account number: %d------------------", i+1))
							b.Send(m.Sender, fmt.Sprintf("Public key(Id): %s\n", sID[i]))
							b.Send(m.Sender, fmt.Sprintf("Secret key(Pw): %s\n", sPW[i]))
						}

						b.Handle(tb.OnText, func(m *tb.Message) {
							v := m.Text
							if s, err := strconv.Atoi(v); err == nil {
								src = sPW[s-1]
								b.Send(m.Sender, "Please enter the receiving account: ")
							}
							b.Handle(tb.OnText, func(m *tb.Message) {
								dst = m.Text
								b.Send(m.Sender, "Please enter the amount to be remitted: ")

								b.Handle(tb.OnText, func(m *tb.Message) {
									amount := m.Text
									b.Send(m.Sender, "Send the amount...")
									resp := stellar.SendPayment(src, dst, amount)
									// Add account balance check func.
									b.Send(m.Sender, "Successful Transaction:")
									b.Send(m.Sender, resp.Ledger)
									b.Send(m.Sender, "check: ")
									b.Send(m.Sender, resp.Hash)
								})
							})
						})

						//stellar.SendPayment()
					} else {

					}
				})
			}

		})

	})

	b.Start()
}
