package main

import (
	"fmt"
	_ "fmt"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}

	address := pair.Address()
	seed := pair.Seed()
	log.Printf("Secret key: %s", seed)
	log.Printf("Public key: %s", address)
	//2021/03/12 19:42:19 Secret key: SAQJCR65OKYXYES2C2Y4YED6Z3PMWKX4ZVIXK6FGSX4M7DJJ753DNMES
	//2021/03/12 19:42:19 Public key: GCKQEDOO6E5BT5OJKHEMJ5NLIMSA3ERDNSJQSC2HJKRVAAKQJK2JEIL3

	//2021/03/12 19:43:06 Secret key: SBP4AWHALPK2MZ2HMOKMVZIZKP4IZVKAYENX322NVQCPJFIOA2ASFOLO
	//2021/03/12 19:43:06 Public key: GBROGZF4YZ5QTPSYEEIQMOWZRZKVDRJUQTJOD5DXNYASM3C4Q7CY553Z

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

}
