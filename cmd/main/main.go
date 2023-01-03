package main

import (
	"flag"

	"github.com/shadmau/MEVProfitExplorer/internal/dashboard"
)

func main() {

	startBlockNumber := flag.Uint("start", 0, "Start Block number")
	endBlockNumber := flag.Uint("end", 0, "End Block number")
	rpc := flag.String("rpc", "", "RPC URL")
	rpcOnly := flag.Bool("rpcOnly", false, "Only use of RPC mode")
	etherscanAPIKey := flag.String("apiKey", "", "Etherscan API Key")
	walletAddress := flag.String("wallet", "", "Wallet address to analyze")

	flag.Parse()
	dashboard.DisplayTextDashboard(*walletAddress, *startBlockNumber, *endBlockNumber, *rpc, *etherscanAPIKey, *rpcOnly)

}
