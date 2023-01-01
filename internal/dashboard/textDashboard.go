package dashboard

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/shadmau/MEVProfitExplorer/pkg/ethereum"
)

const DEFAULT_BLOCKS = 216000 //30 days

func ShowTextDashboard(startBlock uint, endBlock uint, rpc string, ignoreWithdrawls bool) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
		return
	}

	if rpc == "" {
		rpc = os.Getenv("ETH_RPC_URL")
		if rpc == "" {
			log.Fatal("ETH Provider URL not defined in request nor in .env")
		}
	}
	client, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal("Error connecting to Ethereum network:", err)
	}

	currentBlock := ethereum.GetCurrentBlockNumber(client)
	if endBlock == 0 || endBlock > currentBlock {
		endBlock = currentBlock
	}

	if startBlock == 0 || startBlock > currentBlock {
		startBlock = endBlock - DEFAULT_BLOCKS
	}

	//	startBlock = 16306260
	//	endBlock = 16307319
	//	totalProfit := big.NewInt(0)
	for i := startBlock; i < endBlock; i++ {
		/*
				profitForBlock := ethereum.EthProfitForBlock(client, i, "xxx")

				totalProfit.Add(totalProfit, &profitForBlock)
			   fmt.Printf("Analyzing Block %d/%d (%d%%);  Profit: %s ETH \r", i, endBlock, (i-startBlock)*100/(endBlock-startBlock), ethereum.ConvertWEIToETH(totalProfit))
			   time.Sleep(time.Second)
		*/
	}

}

/*

Trading Profit Dashboard

Date: 2022-12-31

Average Profit per Day: 123.45
Average Profit per Week: 678.90
Average Profit per Month: 2987.65

Maximum Profit per Day: 456.78
Maximum Profit per Week: 2345.67
Maximum Profit per Month: 9876.54





*/
