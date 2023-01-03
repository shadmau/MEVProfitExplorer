# MEV Profit Explorer

MEV Profit Explorer is a tool for analyzing the profit made from MEV by a wallet (MEV contract) in a given time period. Currently, this tool is only able to analyze wallets where the profits are made in ETH and stored within the analyzed contract. However, future updates will allow for the analysis of profits made in tokens and expanded support for additional MEV use cases.

## Prerequisites

- Go version 1.15 or higher
- An Ethereum node RPC URL (e.g. Infura or Alchemy)
- An Etherscan API key (optional)

## Installation

1. Clone the repository:
`git clone https://github.com/shadmau/MEVProfitExplorer.git`


2. Install dependencies:
```
go get -u github.com/ethereum/go-ethereum/ethclient
go get -u github.com/joho/godotenv
```

3. (optional) Create a file named `.env` in the project root directory and set the following environment variables:
- `ETH_RPC_URL`: URL of the Ethereum node (e.g. Infura or Alchemy)

## Usage

Run the following command to show the dashboard:
```go run cmd/main/main.go -wallet WALLET_ADDRESS [-start START_BLOCK] [-end END_BLOCK] [-rpc RPC_URL] [-apiKey ETH_SCAN_API_KEY] [-rpcOnly]```



- `WALLET_ADDRESS`: MEV contract address to analyze
- `START_BLOCK`: Block number to start analysis from (optional, default is 30 days before current block)
- `END_BLOCK`: Block number to end analysis at (optional, default is current block)
- `RPC_URL`: URL of the Ethereum node (optional, default value is set as `ETH_RPC_URL` in `.env`)
- `ETH_SCAN_API_KEY`: Etherscan API key (optional)
- `-rpcOnly`: Use only the Ethereum node for data, do not use the Etherscan API (optional)

## Examples

1. `go run ./cmd/main/main.go -wallet 0xCc33Db5fEc8cb1393adD7318Ca99cb916547E1B5  -start 16327123 -end 16344444 -rpc https://eth-mainnet.g.alchemy.com/... `

	This command shows the MEV profit dashboard for the contract `0xCc33Db5fEc8cb1393adD7318Ca99cb916547E1B5` for blocks 16327123 to 16344444 using the specified Alchemy node and  Etherscan API (without API-Key).


2. `go run ./cmd/main/main.go -wallet 0xCc33Db5fEc8cb1393adD7318Ca99cb916547E1B5 `

	This command shows the MEV profit dashboard for the contract `0xCc33Db5fEc8cb1393adD7318Ca99cb916547E1B5` for the last 30 days using an Ethereum node specified in `.env` and  Etherscan API (without API-Key).

#### Example Output
```
Analyzing 109/109 (99%); Total Profit: 0.28766 ETH (+0.00210); AVG Profit/Day: 0.01023 ETH  

MEV Profit Dashboard

From: 2022-12-04 14:12 - 2023-01-03 18:10

Total Profit: 0.28766 ETH
Average Profit per Day: 0.01023 ETH
```
