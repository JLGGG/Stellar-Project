package main

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"log"
)

func main() {
	client := horizonclient.DefaultTestNetClient

	issuerSeed := "SAQJCR65OKYXYES2C2Y4YED6Z3PMWKX4ZVIXK6FGSX4M7DJJ753DNMES"
	distributorSeed := "SBP4AWHALPK2MZ2HMOKMVZIZKP4IZVKAYENX322NVQCPJFIOA2ASFOLO"

	issuer, err := keypair.ParseFull(issuerSeed)
	if err != nil {
		log.Fatal(err)
	}
	distributor, err := keypair.ParseFull(distributorSeed)
	if err != nil {
		log.Fatal(err)
	}

	request := horizonclient.AccountRequest{AccountID: issuer.Address()}
	issuerAccount, err := client.AccountDetail(request)
	if err != nil {
		log.Fatal(err)
	}

	request = horizonclient.AccountRequest{AccountID: distributor.Address()}
	distributorAccount, err := client.AccountDetail(request)
	if err != nil {
		log.Fatal(err)
	}

	// Create an object to represent the new asset
	astroDollar := txnbuild.CreditAsset{Code: "AstroDollar", Issuer: issuer.Address()}

	// First, the receiving(distribution) account must trust the asset from the issuer.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &distributorAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
			Operations: []txnbuild.Operation{
				&txnbuild.ChangeTrust{
					Line:  astroDollar,
					Limit: "5000",
				},
			},
		},
	)

	signedTx, err := tx.Sign(network.TestNetworkPassphrase, distributor)
	resp, err := client.SubmitTransaction(signedTx)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Trust: %s\n", resp.Hash)
	}

	//Second, the issuing account actually sends a payment using the asset
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &issuerAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: distributor.Address(),
					Asset:       astroDollar,
					Amount:      "10",
				},
			},
		})
	signedTx, err = tx.Sign(network.TestNetworkPassphrase, issuer)
	resp, err = client.SubmitTransaction(signedTx)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Pay: %s\n", resp.Hash)
	}
}
