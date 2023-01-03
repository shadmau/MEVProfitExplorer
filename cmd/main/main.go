package main

import (
	"flag"

	"github.com/shadmau/MEVProfitExplorer/internal/dashboard"
)

func main() {

	startBlock := flag.Uint("start", 0, "Start Block number")
	endBlock := flag.Uint("end", 0, "End Block number")
	rpc := flag.String("rpc", "", "RPC")
	rpcOnly := flag.Bool("rpcOnly", false, "RPC")
	ethScanApiKey := flag.String("ethScan", "", "Etherscan API Key")
	walletToAnalyzeStr := flag.String("wallet", "", "Wallet to analyze")

	flag.Parse()
	dashboard.ShowTextDashboard(*walletToAnalyzeStr, *startBlock, *endBlock, *rpc, *ethScanApiKey, *rpcOnly)

}
