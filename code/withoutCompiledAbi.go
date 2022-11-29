package code

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)


type updatePriceData struct {
	description string
	price       float64
	blockNumber uint64
	blockTime   uint64
}


func WithoutCompliledAbi() {

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
	b, err := client.BlockNumber(ctx)
	if err != nil {
		fmt.Println("Failed chainId", err)
	}

	// contract to watch over BTC / USD
	contractAddress := common.HexToAddress("0x4dD6655Ad5ed7C06c882f496E3f42acE5766cb89")
	description := "BTC/USD"
	decimals := int32(8)*-1
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}
	
	logChan := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(ctx, query, logChan)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Opened Listener [%s] on chain ID [%v] block [%v]\n", description, int(chainId.Int64()), big.NewInt(int64(b)))

	transmitAbi, err := abi.JSON(strings.NewReader(TransmitAbi))
	if err != nil {
		fmt.Println("Failed abi", err)
	}

	for {
		select {
			case err := <-sub.Err():
				log.Fatal(err)
			case log := <-logChan:
				if log.Topics[0] == common.HexToHash("0xf6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451") {
					event, err := transmitAbi.EventByID(log.Topics[0])
					if err != nil {
						fmt.Print(err)
					}
					unpacked, err := event.Inputs.Unpack(log.Data)
					if err != nil {
						fmt.Println(err)
					}
					price := unpacked[0].(*big.Int)
					quote := decimal.NewFromBigInt(price, decimals)
					b, err := client.BlockByNumber(ctx, big.NewInt(int64(log.BlockNumber)))
					if err != nil {
						fmt.Println(err)
					}
					d := updatePriceData{
						description: description,
						price: quote.Round(5).InexactFloat64(),
						blockNumber: b.NumberU64(),
						blockTime: b.Time(),
					}
					fmt.Printf("Price Data: %v\n", d)
				}
		}
	}
}