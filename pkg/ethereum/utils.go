package ethereum

import (
	"context"
	"encoding/json"
	"fmt"
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

// EtherScanTransaction represents an Ethereum transaction as returned by the Etherscan API.
type EtherScanTransaction struct {
	BlockNumber string `json:"blockNumber"`
	TimeStamp   string `json:"timeStamp"`
	Hash        string `json:"hash"`
	Nonce       string `json:"nonce"`
	BlockHash   string `json:"blockHash"`
}

// EtherScanResponse represents the response from the Etherscan API.
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

func GetEtherScanTransactionsByAddress(params GetEtherScanTransactionsParams) ([]EtherScanTransaction, error) {
	apiURL := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=%d&endblock=%d&page=1&offset=10000&sort=asc&apikey=%s",
		params.WalletAddress, params.StartBlock, params.EndBlock, params.APIKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	var etherscanResponse EtherScanResponse
	err = json.NewDecoder(resp.Body).Decode(&etherscanResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON etherscan response: %w", err)
	}
	responseStatus, err := strconv.Atoi(etherscanResponse.Status)
	if err != nil {
		return nil, fmt.Errorf("error decoding etherscan response status: %w", err)
	}

	if responseStatus != 1 {
		return nil, fmt.Errorf("etherscan response status is not 1: %d", responseStatus)
	}

	return etherscanResponse.EtherScanTransactions, nil

}

// Returns the timestamp of the given block number.
func BlockToTime(ethClient *ethclient.Client, blocknumber uint) (uint64, error) {
	block, err := ethClient.BlockByNumber(context.Background(), big.NewInt(int64(blocknumber)))
	if err != nil {
		return 0, fmt.Errorf("error fetching block: %w", err)
	}

	return block.Time(), nil
}

// Returns the current block number.
func GetCurrentBlockNumber(ethClient *ethclient.Client) (uint, error) {
	blockNumber, err := ethClient.BlockNumber(context.Background())
	if err != nil {
		return 0, fmt.Errorf("error fetching block: %w", err)
	}
	return uint(blockNumber), nil
}

// Calculates the profit for the given block and wallet address by iterating through all the transactions and checking the wallet balance before and after the block
func MEVProfitForBlock(ethClient *ethclient.Client, blockNumber uint, walletAddress string) (*big.Int, error) {
	blockNumberBefore := int64(blockNumber - 1)
	blockNumberAfter := int64(blockNumber)

	balanceBefore, err := ethClient.BalanceAt(context.Background(), common.HexToAddress(walletAddress), big.NewInt(blockNumberBefore))
	if err != nil {
		return nil, fmt.Errorf("Error fetching balance for block %d: %w", blockNumberBefore, err)
	}
	balanceAfter, err := ethClient.BalanceAt(context.Background(), common.HexToAddress(walletAddress), big.NewInt(blockNumberAfter))
	if err != nil {
		return nil, fmt.Errorf("Error fetching balance for block %d: %w", blockNumberBefore, err)
	}
	profitPerBlock := new(big.Int).Set(balanceAfter).Sub(balanceAfter, balanceBefore)
	if profitPerBlock.Cmp(big.NewInt(0)) > 0 {
		transactionFees, err := GetTransactionFeesByBlock(ethClient, blockNumber, walletAddress)
		if err != nil {
			return nil, fmt.Errorf("error calculating transaction fees for block %d: %w", blockNumber, err)
		}
		profitPerBlock.Sub(profitPerBlock, transactionFees)
	}

	return profitPerBlock, nil
}

// Returns the total transaction fees paid by the given wallet address in the given block. (Fees = Gas Used * Gas Price + Value)
func GetTransactionFeesByBlock(ethClient *ethclient.Client, blocknumber uint, walletAddressString string) (*big.Int, error) {
	// Get the Block
	block, err := ethClient.BlockByNumber(context.Background(), big.NewInt(int64(blocknumber)))
	if err != nil {
		return nil, fmt.Errorf("error fetching block: %w", err)
	}
	totalFeesPaid := big.NewInt(0)
	transactions := block.Transactions()
	for _, transaction := range transactions {
		receiverAddress := transaction.To()
		if receiverAddress != nil && *receiverAddress == common.HexToAddress(walletAddressString) {
			txReceipt, err := ethClient.TransactionReceipt(context.Background(), transaction.Hash())
			if err != nil {
				return nil, fmt.Errorf("error fetching transaction receipt: %w", err)
			}
			gasCosts := big.NewInt(0).Mul(big.NewInt(int64(txReceipt.GasUsed)), transaction.GasPrice())
			txFullGasCosts := new(big.Int).Set(transaction.Value()).Add(transaction.Value(), gasCosts)
			totalFeesPaid.Add(totalFeesPaid, txFullGasCosts)
		}

	}
	return totalFeesPaid, nil
}

// Converts WEI to ETH, with the given precision.
func ConvertWEIToETH(wei *big.Int, precision uint) string {
	weiFloat := new(big.Float).SetInt(wei)
	ethFloat := new(big.Float).Quo(weiFloat, big.NewFloat(1e18))
	eth, _ := ethFloat.Float64()
	return fmt.Sprintf("%.*f", precision, eth)
}
