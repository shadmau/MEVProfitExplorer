package dashboard

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/shadmau/MEVProfitExplorer/pkg/ethereum"
)

const DEFAULT_BLOCKS = 216000 //30 days

func ShowTextDashboard(walletToAnalyzeStr string, startBlock uint, endBlock uint, rpc string, ignoreWithdrawls bool) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	if !common.IsHexAddress(walletToAnalyzeStr) {
		log.Fatal("No valid wallet address specified")
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

	totalProfit := big.NewInt(0)

	for i := startBlock; i < endBlock; i++ {
		profitForBlock := ethereum.EthProfitForBlock(client, i, walletToAnalyzeStr) // if < 0 -> add to Withdrawals list
		totalBlocksAnalyzed := i - startBlock + 1
		avgProfitPerBlock := big.NewInt(0).Div(totalProfit, new(big.Int).SetUint64(uint64(totalBlocksAnalyzed)))
		avgProfitPerDay := new(big.Int).Set(avgProfitPerBlock).Mul(avgProfitPerBlock, big.NewInt(7200)) //7200 blocks per Day
		totalProfit.Add(totalProfit, profitForBlock)
		formatString := "Analyzing %d/%d (%d%%); Total Profit: %s ETH (+%s); AVG Profit/Day: %s ETH  \r"
		fmt.Printf(formatString, i, endBlock, (i-startBlock)*100/(endBlock-startBlock), ethereum.ConvertWEIToETH(totalProfit, 5), ethereum.ConvertWEIToETH(profitForBlock, 5), ethereum.ConvertWEIToETH(avgProfitPerDay, 5))

	}

}

/*

MEV Profit Dashboard

From: xxxx-xx-xx xx:xx - xxxx-xx-xx xx:xx

Total Profit: xx ETH
Average Profit per Day: xx ETH

Withdrawals: xx ETH (xxxx, xxxx, xxxx, xxxx, xxxx, ..)

*/
