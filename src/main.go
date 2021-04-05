package main

import (
	_ "bufio"
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

// Token information is always kept private
const telegramBotToken string = ""

// Save file name
const fileNameAboutEntity string = "src/info.txt"

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func currentDirectory() {
	path, err := os.Getwd()
	checkError(err)
	log.Println(path)
}

func existFile(name string) bool {
	if fi, err := os.Stat(name); err == nil && fi != nil {
		if fi.Mode().IsRegular() {
			return true
		}
	}
	return false
}

func writeFile(fName string, b []byte) {
	err := ioutil.WriteFile(fName, b, os.FileMode(777))
	checkError(err)
}

func writeTxToDB(id, pw string) {
	// Start the SQL server before using SQL: mysql.server start
	// Turn off SQL server: mysql.server stop
	// Keep the sql connection string private.
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/testdb")
	checkError(err)
	defer db.Close()

	result, err := db.Exec("INSERT INTO entity_table (ID, PW) VALUES (?, ?)", id, pw)
	checkError(err)

	n, err := result.RowsAffected()
	log.Printf("%d row inserted\n", n)
}

func readFile(fName string) []byte {
	data, err := ioutil.ReadFile(fName)
	checkError(err)
	return data
}

func readTxFromDB(pID, pPW []string) (count int) {
	// Keep the sql connection string private.
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/testdb")
	checkError(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM entity_table")
	checkError(err)
	defer rows.Close()

	var id, pw string
	count = 0
	for rows.Next() {
		err := rows.Scan(&id, &pw)
		checkError(err)
		pID[count] = id
		pPW[count] = pw
		count += 1
	}
	return
}

func makeAccount() (string, string, string) {
	pair, err := keypair.Random()
	checkError(err)

	address := pair.Address()
	seed := pair.Seed()
	log.Printf("Secret key: %s", seed)
	log.Printf("Public key: %s", address)

	// Get 10,000 test XLM from friendbot.
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + address)
	checkError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	fmt.Println(string(body))

	// Check account information
	request := horizonclient.AccountRequest{AccountID: address}
	account, err := horizonclient.DefaultTestNetClient.AccountDetail(request)
	checkError(err)

	var b bytes.Buffer

	// Used the bytes package for concatenating sentences. It's speed is O(n).
	for _, balance := range account.Balances {
		b.WriteString(fmt.Sprintf("%s\n", balance.Balance))
	}

	return address, seed, b.String()
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
	checkError(err)

	b.Handle("/", func(m *tb.Message) {
		b.Send(m.Sender, "List of Supported Commands:\n")
		b.Send(m.Sender, "/hello")
		b.Send(m.Sender, "/make_account")
		b.Send(m.Sender, "/show_account")
	})

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello World!")
	})

	//TODO What to do when /make_account called twice?
	b.Handle("/make_account", func(m *tb.Message) {
		// keysAndBalance[0]: public key
		// keysAndBalance[1]: private key
		// keysAndBalance[2]: account's balance string
		keysAndBalance[0], keysAndBalance[1], keysAndBalance[2] = makeAccount()
		address := fmt.Sprintf("Public key(Id): %s\n", keysAndBalance[0])
		seed := fmt.Sprintf("Secret key(Pw): %s\n", keysAndBalance[1])

		balanceResult := ParseBalanceStr(keysAndBalance[2])
		buffer.WriteString(fmt.Sprintf("Account ID: https://horizon-testnet.stellar.org/accounts/%s\n", keysAndBalance[0]))
		buffer.WriteString(fmt.Sprintf("Current balance: %s\n", balanceResult))

		//TODO Id, Pw should be saved
		//var temp []byte
		//temp = append(temp, []byte(address)...)
		//temp = append(temp, []byte(seed)...)
		//temp = append(temp, []byte(buffer.String())...)
		//writeFile(fileNameAboutEntity, temp)
		writeTxToDB(keysAndBalance[0], keysAndBalance[1])

		b.Send(m.Sender, address)
		b.Send(m.Sender, seed)
		b.Send(m.Sender, buffer.String())
	})

	b.Handle("/show_account", func(m *tb.Message) {
		//b.Send(m.Sender, string(readFile(fileNameAboutEntity)))
		sID := make([]string, 100)
		sPW := make([]string, 100)
		count := readTxFromDB(sID, sPW)
		for i := 0; i < count; i++ {
			b.Send(m.Sender, sID[i])
			b.Send(m.Sender, sPW[i])
		}
	})

	b.Start()
}
