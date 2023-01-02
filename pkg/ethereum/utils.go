package ethereum

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Converts a Unix timestamp to a block number.
func unixTsToBlock(unixTs uint) (uint, error) {
	// todo: implement
	// https://api.etherscan.io/api?module=block&action=getblocknobytime&timestamp=1669849200&closest=before&apikey=YourApiKeyToken
	return 0, fmt.Errorf("unimplemented")
}

// Returns the timestamp of the given block number.
func BlockToTime(client *ethclient.Client, blocknumber uint) uint64 {
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blocknumber)))
	if err != nil {
		log.Fatal(err)
	}

	return block.Time()
}

// Returns the current block number.
func GetCurrentBlockNumber(client *ethclient.Client) uint {
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal("Could not get Block number!")
	}
	return uint(blockNumber)
}

// Calculates the profit for the given block and wallet address by iterating through all the transactions and checking the wallet balance before and after the block
func EthProfitForBlock(client *ethclient.Client, blockNumber uint, walletAddress string) *big.Int {
	blockNumberBefore := int64(blockNumber - 1)
	blockNumberAfter := int64(blockNumber)

	balanceBefore, err := client.BalanceAt(context.Background(), common.HexToAddress(walletAddress), big.NewInt(blockNumberBefore))
	if err != nil {
		log.Fatal("Could not load balance for " + walletAddress + " for block " + fmt.Sprint(blockNumber))
	}
	balanceAfter, err := client.BalanceAt(context.Background(), common.HexToAddress(walletAddress), big.NewInt(blockNumberAfter))
	if err != nil {
		log.Fatal("Could not load balance for " + walletAddress + " for block " + fmt.Sprint(blockNumber))
	}
	profitPerBlock := new(big.Int).Set(balanceAfter).Sub(balanceAfter, balanceBefore)
	if profitPerBlock.Cmp(big.NewInt(0)) > 0 {
		transactionFees := GetTransactionFeesByBlock(client, blockNumber, walletAddress)
		profitPerBlock.Sub(profitPerBlock, transactionFees)
	}

	return profitPerBlock

}

// Returns the total transaction fees paid by the given wallet address in the given block. (Fees = Gas Used * Gas Price + Value)
func GetTransactionFeesByBlock(client *ethclient.Client, blocknumber uint, walletAddressString string) *big.Int {
	// Get the Block
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blocknumber)))
	if err != nil {
		log.Fatal(err)
	}
	totalFeesPaid := big.NewInt(0)
	transactions := block.Transactions()
	for _, transaction := range transactions {
		receiverAddress := transaction.To()
		if receiverAddress != nil && *receiverAddress == common.HexToAddress(walletAddressString) {
			txReceipt, err := client.TransactionReceipt(context.Background(), transaction.Hash())
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
