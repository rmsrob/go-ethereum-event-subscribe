# Go Ethereum Event Subscribe

> Methods to subscribe on events with go-ethereum

## With Compiled ABI

> Example with [ChainLink Oracle](https://data.chain.link/ethereum/mainnet) compile the ABI
- Copy the abi from Etherscan to the folder `ABIs` in a `.abi` file
- Compile the abi with `abigen` CLI
```sh
abigen --abi=./abis/AccessControlledOffchainAggregator.abi --pkg=abis --out=./abis/AccessControlledOffchainAggregator.go
```

change the smart contract in `main.go`

```go
contractAddress := common.HexToAddress("0xAe74faA92cB67A95ebCAB07358bC222e33A34dA7")
```

## Without Compiled ABI

> Example with [ChainLink Oracle](https://data.chain.link/ethereum/mainnet) in `code` folder
- Need to find the Topics[0] of the event we are looking for
- As we are filter all logs to matche the topics we need
```go
if log.Topics[0] == common.HexToHash("0xf6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451") {
```


update/change the smart contract in `main.go`

- contractAddress
- description
- decimals
```go
contractAddress := common.HexToAddress("0xAe74faA92cB67A95ebCAB07358bC222e33A34dA7")
description := "BTC/USD"
decimals := int32(8)*-1
```



## Usage

```sh
go run main.go
```
