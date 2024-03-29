package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
	"github.com/opefago/bot-tave/exchange/base"
	"github.com/opefago/bot-tave/exchange/poloniex"
	"github.com/opefago/bot-tave/utils"
)

var (
	NodeEndpoint     = ""
	SampleOutputFile = "Poloniex.json"
	BinanceBaseUrl   = ""
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	NodeEndpoint = os.Getenv("NODE_ENDPOINT")
}

func main() {
	client, err := ethclient.Dial(NodeEndpoint)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		client.Close()
	}()

	cntxt := context.Background()
	txnsHash := make(chan common.Hash)

	var exchange base.Exchange
	exchange = poloniex.New()

	// err = exchange.StoreTriangularPair()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// go ListenMemPool(cntxt, client, txnsHash)

	exchange.RunExchange(cntxt)
	// fmt.Println()
	// fmt.Println("======================================================")
	// fmt.Println()

	// exchange = binance.New()

	// err = exchange.StoreTriangularPair()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	_ = cntxt
	_ = txnsHash

	for {
	}
}

func ListenMemPool(ctx context.Context, client *ethclient.Client, channel chan common.Hash) error {
	baseClient, err := rpc.Dial(NodeEndpoint)
	if err != nil {
		return err
	}

	fmt.Println("RPC connection successful!")

	defer func() {
		baseClient.Close()
	}()

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return err
	}

	subscriber := gethclient.New(baseClient)
	_, err = subscriber.SubscribePendingTransactions(ctx, channel)

	if err != nil {
		return err
	}

	fmt.Println("Pending subscription connected successful!")

	signer := types.NewLondonSigner(chainID)

	for txnHash := range channel {
		txn, _, err := client.TransactionByHash(ctx, txnHash)
		if err != nil {
			continue
		}

		message, err := txn.AsMessage(signer, nil)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println()
		fmt.Println("-----------------------------------------")
		fmt.Println()
		fmt.Printf("Transaction Hash: %s\n", txnHash.String())
		fmt.Printf("To Address: %s\n", message.To())
		fmt.Printf("Value: %v\n", utils.WeiToETH(txn.Value()))
		fmt.Printf("Nonce: %d\n", txn.Nonce())
		fmt.Printf("Gas: %d\n", txn.Gas())
		fmt.Printf("Data: %s\n", txn.Data())

		fmt.Println()
		fmt.Println("-----------------------------------------")
		fmt.Println()
	}

	return nil
}
