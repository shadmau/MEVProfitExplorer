package dashboard

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/shadmau/MEVProfitExplorer/pkg/ethereum"
)

const DEFAULT_BLOCKS = 216000 //30 days

/* Output:
MEV Profit Dashboard

Time: YYYY-MM-DD HH:SS - YYYY-MM-DD HH:SS

Total Profit: 1.00000 ETH
Average Profit per Day: 0.10000 ETH

Withdrawals: 1.20000 ETH (16314680,16314681,16314690)
*/

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
	var withdrawals []uint64

	totalProfit := big.NewInt(0)
	avgProfitPerDay := big.NewInt(0)
	totalWithdrawals := big.NewInt(0)

	for i := startBlock; i < endBlock; i++ {
		profitForBlock := ethereum.EthProfitForBlock(client, i, walletToAnalyzeStr)
		if profitForBlock.Cmp(big.NewInt(0)) >= 0 {
			totalProfit.Add(totalProfit, profitForBlock)
		} else {
			totalWithdrawals.Add(totalWithdrawals, profitForBlock)
			withdrawals = append(withdrawals, uint64(i))
		}
		totalBlocksAnalyzed := i - startBlock + 1
		avgProfitPerBlock := big.NewInt(0).Div(totalProfit, new(big.Int).SetUint64(uint64(totalBlocksAnalyzed)))
		avgProfitPerDay = new(big.Int).Set(avgProfitPerBlock).Mul(avgProfitPerBlock, big.NewInt(7200)) //7200 blocks per Day
		formatString := "Analyzing %d/%d (%d%%); Total Profit: %s ETH (+%s); AVG Profit/Day: %s ETH  \r"
		fmt.Printf(formatString, i, endBlock, (i-startBlock)*100/(endBlock-startBlock), ethereum.ConvertWEIToETH(totalProfit, 5), ethereum.ConvertWEIToETH(profitForBlock, 5), ethereum.ConvertWEIToETH(avgProfitPerDay, 5))

	}
	fmt.Println()

	startTime := time.Unix(int64(ethereum.BlockToTime(client, startBlock)), 0)
	startTimeString := startTime.Format("2006-01-02 15:04")

	endTime := time.Unix(int64(ethereum.BlockToTime(client, endBlock)), 0)
	endTimeString := endTime.Format("2006-01-02 15:04")
	withdrawalsString := ethereum.ConvertWEIToETH(totalWithdrawals, 5) + " ETH "
	if len(withdrawals) > 0 {
		withdrawalsString += "("
		for i, v := range withdrawals {
			withdrawalsString += fmt.Sprint(v)
			if i+1 < len(withdrawals) {
				withdrawalsString += ","
			}
		}
		withdrawalsString += ")"
	}
	averageProfitString := "analyze period too short"
	if endBlock-startBlock >= 7200 {
		averageProfitString = ethereum.ConvertWEIToETH(avgProfitPerDay, 5) + " ETH"
	}
	fmt.Printf("MEV Profit Dashboard\n\nFrom: %s - %s\n\nTotal Profit: %s ETH\nAverage Profit per Day: %s\n\nWithdrawals: %s \n", startTimeString, endTimeString, ethereum.ConvertWEIToETH(totalProfit, 5), averageProfitString, withdrawalsString)

}
