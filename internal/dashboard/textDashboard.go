package dashboard

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
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

func ShowTextDashboard(walletToAnalyze string, startBlock, endBlock uint, rpc, apiKey string, rpcOnly bool) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	if !common.IsHexAddress(walletToAnalyze) {
		log.Fatal("No valid wallet address specified")
	}

	if rpc == "" {
		rpc = os.Getenv("ETH_RPC_URL")
		if rpc == "" {
			log.Fatal("ETH Provider URL not defined in request nor in .env")
		}
	}
	ethClient, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal("Error connecting to Ethereum network:", err)
	}

	currentBlock := ethereum.GetCurrentBlockNumber(ethClient)
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
	var allBlockNumbers []uint
	if rpcOnly {
		for i := startBlock; i < endBlock; i++ {
			allBlockNumbers = append(allBlockNumbers, i)
		}

	} else {
		params := ethereum.GetEtherScanTransactionsParams{
			APIKey:        apiKey,
			WalletAddress: walletToAnalyze,
			StartBlock:    uint64(startBlock),
			EndBlock:      uint64(endBlock),
		}
		transactions := ethereum.GetEtherScanTransactionsByAddress(params)
		for _, transaction := range transactions {
			blockNumber, _ := strconv.Atoi(transaction.BlockNumber)
			allBlockNumbers = append(allBlockNumbers, uint(blockNumber))
		}
	}

	for i, blockNumber := range allBlockNumbers {
		profitForBlock := ethereum.EthProfitForBlock(ethClient, blockNumber, walletToAnalyze)
		if profitForBlock.Cmp(big.NewInt(0)) >= 0 {
			totalProfit.Add(totalProfit, profitForBlock)
		} else {
			totalWithdrawals.Add(totalWithdrawals, profitForBlock)
			withdrawals = append(withdrawals, uint64(blockNumber))
		}
		blockNumberAnalyzed := i + 1
		totalBlocksAnalyzed := big.NewInt(int64(allBlockNumbers[i])).Sub(big.NewInt(int64(allBlockNumbers[i])+1), big.NewInt(int64(allBlockNumbers[0])))
		avgProfitPerBlock := big.NewInt(0).Div(totalProfit, totalBlocksAnalyzed)
		avgProfitPerDay = new(big.Int).Set(avgProfitPerBlock).Mul(avgProfitPerBlock, big.NewInt(7200)) //7200 blocks per Day
		formatString := "Analyzing %d/%d (%d%%); Total Profit: %s ETH (+%s); AVG Profit/Day: %s ETH  \r"
		fmt.Printf(formatString, blockNumberAnalyzed, len(allBlockNumbers), i*100/len(allBlockNumbers), ethereum.ConvertWEIToETH(totalProfit, 5), ethereum.ConvertWEIToETH(profitForBlock, 5), ethereum.ConvertWEIToETH(avgProfitPerDay, 5))

	}
	fmt.Println("\r\r\n")

	startTime := time.Unix(int64(ethereum.BlockToTime(ethClient, startBlock)), 0)
	startTimeString := startTime.Format("2006-01-02 15:04")

	endTime := time.Unix(int64(ethereum.BlockToTime(ethClient, endBlock)), 0)
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
