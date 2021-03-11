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
	//2021/03/10 07:28:45 Secret key: SBCRHU24LEAUGZ7PMFAXNPCZM34JGGYD2N3CPILQIGQWBZFNS7VOCDZ6
	//2021/03/10 07:28:45 Public key: GA7KLAUPA5RTQQ56NN2S3PV5A2FTBX2MDVNMC6SI74FJSO4GGGSAY2TN

	//2021/03/10 07:29:38 Secret key: SAWRQSCT3BLGVSBWZCUB76M4CCSNA6GME4X5JPC36UD6524BBJBJLGCE
	//2021/03/10 07:29:38 Public key: GDYO22NN6UCRI2TRPFAR4RC4LQFFZU673RF4LTZBFUJDUECG5QVMINHM

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
