package stellar

import (
	"bytes"
	"fmt"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func CurrentDirectory() {
	path, err := os.Getwd()
	CheckError(err)
	log.Println(path)
}

func ExistFile(name string) bool {
	if fi, err := os.Stat(name); err == nil && fi != nil {
		if fi.Mode().IsRegular() {
			return true
		}
	}
	return false
}

func WriteFile(fName string, b []byte) {
	err := ioutil.WriteFile(fName, b, os.FileMode(777))
	CheckError(err)
}

func ReadFile(fName string) []byte {
	data, err := ioutil.ReadFile(fName)
	CheckError(err)
	return data
}

func MakeAccount() (string, string, string) {
	pair, err := keypair.Random()
	CheckError(err)

	address := pair.Address()
	seed := pair.Seed()
	log.Printf("Public key: %s", address)
	log.Printf("Secret key: %s", seed)

	// Get 10,000 test XLM from friendbot.
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + address)
	CheckError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	CheckError(err)
	fmt.Println(string(body))

	// Check account information
	request := horizonclient.AccountRequest{AccountID: address}
	account, err := horizonclient.DefaultTestNetClient.AccountDetail(request)
	CheckError(err)

	var b bytes.Buffer

	// Used the bytes package for concatenating sentences. It's speed is O(n).
	for _, balance := range account.Balances {
		b.WriteString(fmt.Sprintf("%s\n", balance.Balance))
	}

	return address, seed, b.String()
}

func SendPayment(src, dest, amount string) horizon.Transaction {
	client := horizonclient.DefaultTestNetClient

	// Make sure destination account exists
	destAccountRequest := horizonclient.AccountRequest{AccountID: dest}
	_, err := client.AccountDetail(destAccountRequest)
	CheckError(err)

	// Load the source account
	sourceKP := keypair.MustParseFull(src)
	sourceAccountRequest := horizonclient.AccountRequest{AccountID: sourceKP.Address()}
	sourceAccount, err := client.AccountDetail(sourceAccountRequest)
	CheckError(err)

	// Build transaction
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(), // Use a real timeout in production!
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: dest,
					Amount:      amount,
					Asset:       txnbuild.NativeAsset{},
				},
			},
		},
	)
	CheckError(err)

	// Sign the transaction to prove you are actually the person sending it.
	tx, err = tx.Sign(network.TestNetworkPassphrase, sourceKP)
	CheckError(err)

	// And finally, send it off to Stellar!
	resp, err := horizonclient.DefaultTestNetClient.SubmitTransaction(tx)
	CheckError(err)

	return resp
}

func ParseBalanceStr(balanceStr string) string {
	// need to modify regular expression.
	regexp := regexp.MustCompile("[0-9]+\\.[0-9]+")
	balanceStr = regexp.FindAllString(balanceStr, 1)[0]

	return balanceStr
}

func CheckAccountBalance(balance, requestedAmount string) (float64, bool) {
	s, _ := strconv.ParseFloat(balance, 64)
	v, _ := strconv.ParseFloat(requestedAmount, 64) // requested amount
	if s < v {
		return s, false // Remittance amount exceeds the balance in the account.
	} else {
		return s, true
	}
}

func ReturnBalance(address string) string {
	client := horizonclient.DefaultTestNetClient

	// Make sure destination account exists
	request := horizonclient.AccountRequest{AccountID: address}
	account, err := client.AccountDetail(request)
	CheckError(err)

	var b bytes.Buffer

	// Used the bytes package for concatenating sentences. It's speed is O(n).
	for _, balance := range account.Balances {
		b.WriteString(fmt.Sprintf("%s\n", balance.Balance))
	}

	return ParseBalanceStr(b.String())
}
