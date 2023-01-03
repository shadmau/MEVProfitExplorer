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

Withdrawals: 1.20000 ETH (16314680, 16314681, 16314690)

*/

// Dissplays a text-based dashboard with information about MEV profits.
func DisplayTextDashboard(walletAddress string, startBlockNumber, endBlockNumber uint, rpcURL, etherscanAPIKey string, rpcOnlyMode bool) {

	// Load environment variables from .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Check that a valid wallet address was specified
	if !common.IsHexAddress(walletAddress) {
		log.Fatal("No valid wallet address specified")
	}

	// Use env RPC URL if none was specified
	if rpcURL == "" {
		rpcURL = os.Getenv("ETH_RPC_URL")
		if rpcURL == "" {
			log.Fatal("ETH Provider URL not defined in request nor in .env")
		}
	}

	// Connect to Ethereum network
	ethClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("Error connecting to Ethereum network:", err)
	}

	// Get current block number
	currentBlock, err := ethereum.GetCurrentBlockNumber(ethClient)
	if err != nil {
		log.Fatal("error fetching current blocknumber: %w", err)
	}

	// Use if no end block specified us current block as end block
	if endBlockNumber == 0 || endBlockNumber > currentBlock {
		endBlockNumber = currentBlock
	}

	// Use if no start block specified use a block approx. 30 days ago
	if startBlockNumber == 0 || startBlockNumber > currentBlock {
		startBlockNumber = endBlockNumber - DEFAULT_BLOCKS
	}
	var withdrawals []uint64

	totalProfit := big.NewInt(0)
	avgProfitPerDay := big.NewInt(0)
	totalWithdrawals := big.NewInt(0)
	var allBlockNumbers []uint
	if rpcOnlyMode {
		// Add all block numbers to allBlockNumbers
		for i := startBlockNumber; i < endBlockNumber; i++ {
			allBlockNumbers = append(allBlockNumbers, i)
		}

	} else {

		etherscanParams := ethereum.GetEtherScanTransactionsParams{
			APIKey:        etherscanAPIKey,
			WalletAddress: walletAddress,
			StartBlock:    uint64(startBlockNumber),
			EndBlock:      uint64(endBlockNumber),
		}
		// Get all transactions that have been sent to the walletAddress
		transactions, err := ethereum.GetEtherScanTransactionsByAddress(etherscanParams)
		if err != nil {
			log.Fatal("error fetching etherscan transactions: %w", err)
		}

		// Add only block numbers where a transaction has been sent to the wallet to allBlockNumbers
		for _, transaction := range transactions {
			blockNumber, _ := strconv.Atoi(transaction.BlockNumber)
			allBlockNumbers = append(allBlockNumbers, uint(blockNumber))
		}
	}

	for i, blockNumber := range allBlockNumbers {
		//Get MEV Profit for wallet in block number
		profitForBlock, err := ethereum.MEVProfitForBlock(ethClient, blockNumber, walletAddress)
		if err != nil {
			log.Fatal("error fetching profit for block: %w", err)
		}
		if profitForBlock.Cmp(big.NewInt(0)) >= 0 {
			totalProfit.Add(totalProfit, profitForBlock)
		} else {
			// Every transaction that removed funds from wallet is considered as a withdrawl
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

	fmt.Println("")

	startTime, err := ethereum.BlockToTime(ethClient, startBlockNumber)
	if err != nil {
		startTime = 0
		fmt.Printf("error fetching block time for block %d: %v", startBlockNumber, err)
	}

	startTimeUnix := time.Unix(int64(startTime), 0)
	startTimeString := startTimeUnix.Format("2006-01-02 15:04")
	endTime, err := ethereum.BlockToTime(ethClient, endBlockNumber)
	if err != nil {
		endTime = 0
		fmt.Printf("error fetching block time for block %d: %v", endBlockNumber, err)
	}
	endTimeUnix := time.Unix(int64(endTime), 0)
	endTimeString := endTimeUnix.Format("2006-01-02 15:04")

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

	// Profit per day will only be shown if at least 7200 blocks (1 day) have been analyzed
	if endBlockNumber-startBlockNumber >= 7200 {
		averageProfitString = ethereum.ConvertWEIToETH(avgProfitPerDay, 5) + " ETH"
	}

	fmt.Printf("\r\nMEV Profit Dashboard\t\n\nFrom: %s - %s\n\nTotal Profit: %s ETH\nAverage Profit per Day: %s\n\nWithdrawals: %s \n", startTimeString, endTimeString, ethereum.ConvertWEIToETH(totalProfit, 5), averageProfitString, withdrawalsString)

}
