package main

import (
	"flag"

	"github.com/shadmau/MEVProfitExplorer/internal/dashboard"
)

func main() {

	//startTs := flag.Int("startTs", 0, "Start Time to analyze")
	//endTs := flag.Int("endTs", 0, "End of analyze")
	startBlock := flag.Uint("start", 0, "Start Block number")
	endBlock := flag.Uint("end", 0, "End Block number")
	rpc := flag.String("rpc", "", "RPC")
	ignoreWithdrawls := flag.Bool("igno", true, "Ignore withdrawals")
	flag.Parse()

	dashboard.ShowTextDashboard(*startBlock, *endBlock, *rpc, *ignoreWithdrawls)

}
