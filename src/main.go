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

// Save file name
const fileNameAboutEntity string = "src/info.txt"

func main() {
	// TODO:
	//	   Add account favorites

	// Token information is always kept private
	// Input your bot token
	//if len(os.Args) < 2 {
	//	panic("Error: less than two arguments")
	//}
	//telegramBotToken := os.Args[1]
	telegramBotToken := ""

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
		b.Send(m.Sender, "/show_favorite")
		b.Send(m.Sender, "/save_favorite")
		b.Send(m.Sender, "/delete_favorite")
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

		stellar.WriteTxToDB(keysAndBalance[0], keysAndBalance[1])

		b.Send(m.Sender, address)
		b.Send(m.Sender, seed)
		b.Send(m.Sender, buffer.String())
		buffer.Reset()
	})

	b.Handle("/show_account", func(m *tb.Message) {
		sID := make([]string, 100)
		sPW := make([]string, 100)
		count := stellar.ReadTxFromDB(sID, sPW)

		b.Send(m.Sender, "View all your account information")
		for i := 0; i < count; i++ {
			strID := fmt.Sprintf("Public key(Id): %s\n", sID[i])
			strPW := fmt.Sprintf("Secret key(Pw): %s\n", sPW[i])
			note := fmt.Sprintf("------------------ Account number: %d------------------", i+1)
			b.Send(m.Sender, note)
			b.Send(m.Sender, strID)
			b.Send(m.Sender, strPW)
		}
	})

	b.Handle("/show_favorite", func(m *tb.Message) {
		b.Send(m.Sender, "This command displays a list of favorite accounts.")

		sID := make([]string, 100)
		count := stellar.ReadFavoriteAccountFromDB(sID)

		b.Send(m.Sender, "View all favorite account information")
		for i := 0; i < count; i++ {
			b.Send(m.Sender, fmt.Sprintf("------------------ Account number: %d------------------", i+1))
			b.Send(m.Sender, fmt.Sprintf("Public key(Id): %s\n", sID[i]))
		}
	})

	b.Handle("/save_favorite", func(m *tb.Message) {
		b.Send(m.Sender, "This command saves a frequently used account.")
		b.Send(m.Sender, "Please enter receiving account to be saved.")
		b.Handle(tb.OnText, func(m *tb.Message) {
			str := m.Text
			stellar.WriteFavoriteAccountToDB(str)
		})

		b.Send(m.Sender, "Saved")
		sID := make([]string, 100)
		count := stellar.ReadFavoriteAccountFromDB(sID)

		b.Send(m.Sender, "View all favorite account information")
		for i := 0; i < count; i++ {
			b.Send(m.Sender, fmt.Sprintf("------------------ Account number: %d------------------", i+1))
			b.Send(m.Sender, fmt.Sprintf("Public key(Id): %s\n", sID[i]))
		}

	})

	b.Handle("/delete_favorite", func(m *tb.Message) {

	})

	b.Handle("/send_payment", func(m *tb.Message) {
		// Enter the address to send
		// Check frequently used accounts
		b.Send(m.Sender, fmt.Sprintln("This is a remittance command."))
		b.Send(m.Sender, fmt.Sprint("Would you like to browse your frequently used accounts? (y/n): "))
		b.Handle(tb.OnText, func(m *tb.Message) {
			var srcAddress, srcSeed, dst, balance string
			sID := make([]string, 100)
			sPW := make([]string, 100)
			dID := make([]string, 100)

			if m.Text == "y" {
				count := stellar.ReadFavoriteAccountFromDB(dID)
				b.Send(m.Sender, "The list of favorite accounts:")
				b.Send(m.Sender, "Select the receiving account to use when sending money")
				for i := 0; i < count; i++ {
					b.Send(m.Sender, fmt.Sprintf("------------------ Account number: %d------------------", i+1))
					b.Send(m.Sender, fmt.Sprintf("Public key(Id): %s\n", dID[i]))
				}

				b.Handle(tb.OnText, func(m *tb.Message) {
					v := m.Text
					if s, err := strconv.Atoi(v); err == nil {
						dst = dID[s-1]
						b.Send(m.Sender, fmt.Sprint("Enter the address to send(Press \"list\" to see the list of your accounts):"))

						b.Handle(tb.OnText, func(m *tb.Message) {
							if len(m.Text) != 56 && m.Text != "list" {
								b.Send(m.Sender, "Account ID is 56 characters. Please re-enter.")
							} else if m.Text == "list" {
								b.Send(m.Sender, "Please select the remittance account you want to use:")
								count := stellar.ReadTxFromDB(sID, sPW)
								for i := 0; i < count; i++ {
									b.Send(m.Sender, fmt.Sprintf("------------------ Account number: %d------------------", i+1))
									b.Send(m.Sender, fmt.Sprintf("Public key(Id): %s\n", sID[i]))
									b.Send(m.Sender, fmt.Sprintf("Secret key(Pw): %s\n", sPW[i]))
								}
								b.Handle(tb.OnText, func(m *tb.Message) {
									v := m.Text
									if s, err := strconv.Atoi(v); err == nil {
										srcSeed = sPW[s-1]
										srcAddress = sID[s-1]
									}
									b.Handle(tb.OnText, func(m *tb.Message) {
										balance = stellar.ReturnBalance(srcAddress)
										b.Send(m.Sender, fmt.Sprintf("Your current balance: %s", balance))
										b.Send(m.Sender, "Please enter the amount to be remitted: ")

										b.Handle(tb.OnText, func(m *tb.Message) {
											amount := m.Text
											if s, err := stellar.CheckAccountBalance(balance, m.Text); err == false {
												b.Send(m.Sender, fmt.Sprintf("The remittance amount exceeds the account's balance. Your current balance: %f. Please re-enter.", s))
												time.Sleep(10 * time.Second)
											}
											b.Send(m.Sender, "Send the amount...")
											resp := stellar.SendPayment(srcSeed, dst, amount)
											// Add account balance check func.
											b.Send(m.Sender, "Successful Transaction:")
											b.Send(m.Sender, fmt.Sprintf("https://horizon-testnet.stellar.org/accounts/%s", dst))
											b.Send(m.Sender, fmt.Sprintf("Check: %s", resp.Hash))
										})
									})
								})
							}
						})
					}
				})
			} else {
				var srcAddress, srcSeed, dst, balance string
				sID := make([]string, 100)
				sPW := make([]string, 100)
				b.Send(m.Sender, fmt.Sprint("Enter the address to send(Press \"list\" to see the list of your accounts):"))

				b.Handle(tb.OnText, func(m *tb.Message) {
					if len(m.Text) != 56 && m.Text != "list" {
						b.Send(m.Sender, "Account ID is 56 characters. Please re-enter.")
					} else if m.Text == "list" {
						b.Send(m.Sender, "Please select the remittance account you want to use:")

						count := stellar.ReadTxFromDB(sID, sPW)
						for i := 0; i < count; i++ {
							b.Send(m.Sender, fmt.Sprintf("------------------ Account number: %d------------------", i+1))
							b.Send(m.Sender, fmt.Sprintf("Public key(Id): %s\n", sID[i]))
							b.Send(m.Sender, fmt.Sprintf("Secret key(Pw): %s\n", sPW[i]))
						}

						b.Handle(tb.OnText, func(m *tb.Message) {
							v := m.Text
							if s, err := strconv.Atoi(v); err == nil {
								srcSeed = sPW[s-1]
								srcAddress = sID[s-1]
								b.Send(m.Sender, "Please enter the receiving account: ")
							}

							b.Handle(tb.OnText, func(m *tb.Message) {
								// 송금시 입력한 수신계좌 즐겨찾기 추가 코드
								dst = m.Text
								balance = stellar.ReturnBalance(srcAddress)
								b.Send(m.Sender, fmt.Sprintf("Your current balance: %s", balance))
								b.Send(m.Sender, "Please enter the amount to be remitted: ")

								b.Handle(tb.OnText, func(m *tb.Message) {
									amount := m.Text
									if s, err := stellar.CheckAccountBalance(balance, m.Text); err == false {
										b.Send(m.Sender, fmt.Sprintf("The remittance amount exceeds the account's balance. Your current balance: %f. Please re-enter.", s))
										time.Sleep(10 * time.Second)
									}
									b.Send(m.Sender, "Send the amount...")
									resp := stellar.SendPayment(srcSeed, dst, amount)
									// Add account balance check func.
									b.Send(m.Sender, "Successful Transaction:")
									b.Send(m.Sender, fmt.Sprintf("https://horizon-testnet.stellar.org/accounts/%s", dst))
									b.Send(m.Sender, fmt.Sprintf("Check: %s", resp.Hash))
								})
							})
						})
					}
				})
			}

		})

	})

	b.Start()
}
