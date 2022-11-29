package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"

	"github.com/rrobrms/go-ethereum-event-subscribe/abis"
)


type updatePriceData struct {
	description string
	price       float64
	blockNumber uint64
	blockTime   uint64
}


func main() {

	client, err := ethclient.Dial("wss://eth-mainnet.g.alchemy.com/v2/")
	if err != nil {
		fmt.Println("Failed client", err)
	}
	defer client.Close()

	ctx := context.Background()
	chainId, err := client.ChainID(ctx)
	if err != nil {
		fmt.Println("Failed chainId", err)
	}
	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		fmt.Println("Failed chainId", err)
	}

	// contract to watch over BTC / USD
	contractAddress := common.HexToAddress("0xAe74faA92cB67A95ebCAB07358bC222e33A34dA7")

	// created a instance of the Compiled abi
    instance, err := abis.NewAbis(contractAddress, client)
    if err != nil {
        log.Fatal(err)
    }
    description, err := instance.Description(&bind.CallOpts{})
    if err != nil {
        log.Fatal(err)
    }
	decimals, err := instance.Decimals(&bind.CallOpts{})
    if err != nil {
        log.Fatal(err)
    }
	fmt.Printf("Opened Listener [%s] on chain ID [%v] -- block [%v]\n", description, int(chainId.Int64()), big.NewInt(int64(blockNumber)))

	// create channel to accept event types "Transmission" from compiled abi
	wntChan := make(chan *abis.AbisNewTransmission)
	e, err := instance.WatchNewTransmission(&bind.WatchOpts{},wntChan, []uint32{})
	if err != nil {
		fmt.Print(err)
	}	

	for {
		select {
			case err := <-e.Err():
				log.Fatal(err)
			case wt := <-wntChan:
				answer := decimal.NewFromBigInt(wt.Answer, int32(decimals)*-1)
				b, err := client.BlockByNumber(ctx, nil)
				if err != nil {
					fmt.Println(err)
				}
				d := updatePriceData{
					description: description,
					price: answer.Round(5).InexactFloat64(),
					blockNumber: b.NumberU64(),
					blockTime: b.Time(),
				}
				fmt.Printf("Update Price Data: %v\n", d)
		}
	}
}