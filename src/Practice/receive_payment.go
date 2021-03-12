package main

import (
	"context"
	"fmt"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon/operations"
	"time"
)

func main() {
	client := horizonclient.DefaultTestNetClient
	// Account you want to check
	opRequest := horizonclient.OperationRequest{ForAccount: "GBM37UDJOQBEULCPB6QVNTHCIT572VDIG347G5ZYXZYB5PVEREQCPHU6", Cursor: "now"}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Stop streaming after 60 seconds.
		time.Sleep(60 * time.Second)
		cancel()
	}()

	printHandler := func(op operations.Operation) {
		fmt.Println(op)
	}
	err := client.StreamPayments(ctx, opRequest, printHandler)
	if err != nil {
		fmt.Println(err)
	}

}
