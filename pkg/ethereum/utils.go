package ethereum

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Converts a Unix timestamp to a block number.
func unixTsToBlock(unixTs uint) (uint, error) {
	// todo: implement
	// https://api.etherscan.io/api?module=block&action=getblocknobytime&timestamp=1669849200&closest=before&apikey=YourApiKeyToken
	return 0, fmt.Errorf("unimplemented")
}

// Transaction gathered through etherscan API
type EtherScanTransaction struct {
	BlockNumber string `json:"blockNumber"`
	TimeStamp   string `json:"timeStamp"`
	Hash        string `json:"hash"`
	Nonce       string `json:"nonce"`
	BlockHash   string `json:"blockHash"`
}

type EtherScanResponse struct {
	Status                string                 `json:"status"`
	Message               string                 `json:"message"`
	EtherScanTransactions []EtherScanTransaction `json:"result"`
}

type GetEtherScanTransactionsParams struct {
	APIKey        string
	WalletAddress string
	StartBlock    uint64
	EndBlock      uint64
}

func GetEtherScanTransactionsByAddress(params GetEtherScanTransactionsParams) []EtherScanTransaction {
	apiURL := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=%d&endblock=%d&page=1&offset=10000&sort=asc&apikey=%s",
		params.WalletAddress, params.StartBlock, params.EndBlock, params.APIKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var etherscanResponse EtherScanResponse
	err = json.NewDecoder(resp.Body).Decode(&etherscanResponse)
	if err != nil {
		log.Fatal(err)
	}
	responseStatus, err := strconv.Atoi(etherscanResponse.Status)
	if err != nil {
		log.Fatal(err)
	}

	if responseStatus != 1 {
		log.Fatal(etherscanResponse.Message)

	}
	return etherscanResponse.EtherScanTransactions

}

// Returns the timestamp of the given block number.
func BlockToTime(ethClient *ethclient.Client, blocknumber uint) uint64 {
	block, err := ethClient.BlockByNumber(context.Background(), big.NewInt(int64(blocknumber)))
	if err != nil {
		log.Fatal(err)
	}

	return block.Time()
}

// Returns the current block number.
func GetCurrentBlockNumber(ethClient *ethclient.Client) uint {
	blockNumber, err := ethClient.BlockNumber(context.Background())
	if err != nil {
		log.Fatal("Could not get Block number!")
	}
	return uint(blockNumber)
}

// Calculates the profit for the given block and wallet address by iterating through all the transactions and checking the wallet balance before and after the block
func EthProfitForBlock(ethClient *ethclient.Client, blockNumber uint, walletAddress string) *big.Int {
	blockNumberBefore := int64(blockNumber - 1)
	blockNumberAfter := int64(blockNumber)

	balanceBefore, err := ethClient.BalanceAt(context.Background(), common.HexToAddress(walletAddress), big.NewInt(blockNumberBefore))
	if err != nil {
		log.Fatal("Could not load balance for " + walletAddress + " for block " + fmt.Sprint(blockNumber))
	}
	balanceAfter, err := ethClient.BalanceAt(context.Background(), common.HexToAddress(walletAddress), big.NewInt(blockNumberAfter))
	if err != nil {
		log.Fatal("Could not load balance for " + walletAddress + " for block " + fmt.Sprint(blockNumber))
	}
	profitPerBlock := new(big.Int).Set(balanceAfter).Sub(balanceAfter, balanceBefore)
	if profitPerBlock.Cmp(big.NewInt(0)) > 0 {
		transactionFees := GetTransactionFeesByBlock(ethClient, blockNumber, walletAddress)
		profitPerBlock.Sub(profitPerBlock, transactionFees)
	}

	return profitPerBlock
}

// Returns Transactionhashes, Blocknumber

// Returns the total transaction fees paid by the given wallet address in the given block. (Fees = Gas Used * Gas Price + Value)
func GetTransactionFeesByBlock(ethClient *ethclient.Client, blocknumber uint, walletAddressString string) *big.Int {
	// Get the Block
	block, err := ethClient.BlockByNumber(context.Background(), big.NewInt(int64(blocknumber)))
	if err != nil {
		log.Fatal(err)
	}
	totalFeesPaid := big.NewInt(0)
	transactions := block.Transactions()
	for _, transaction := range transactions {
		receiverAddress := transaction.To()
		if receiverAddress != nil && *receiverAddress == common.HexToAddress(walletAddressString) {
			txReceipt, err := ethClient.TransactionReceipt(context.Background(), transaction.Hash())
			if err != nil {
				log.Fatal(err)
			}
			gasCosts := big.NewInt(0).Mul(big.NewInt(int64(txReceipt.GasUsed)), transaction.GasPrice())
			txFullGasCosts := new(big.Int).Set(transaction.Value()).Add(transaction.Value(), gasCosts)
			totalFeesPaid.Add(totalFeesPaid, txFullGasCosts)
		}

	}
	return totalFeesPaid
}

// Converts WEI to ETH, with the given precision.
func ConvertWEIToETH(wei *big.Int, precision uint) string {
	weiFloat := new(big.Float).SetInt(wei)
	ethFloat := new(big.Float).Quo(weiFloat, big.NewFloat(1e18))
	eth, _ := ethFloat.Float64()
	return fmt.Sprintf("%.*f", precision, eth)
}
