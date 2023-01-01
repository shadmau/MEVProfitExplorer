package ethereum

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// todo: implement
func unixTsToBlock(unixTs uint) uint {
	return 0
}

func BlockToTime(client *ethclient.Client, blocknumber uint) uint64 {
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blocknumber)))
	if err != nil {
		log.Fatal(err)
	}

	return block.Time()
}

func GetCurrentBlockNumber(client *ethclient.Client) uint {
	blocknumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal("Could not get Blocknumber!")
	}
	return uint(blocknumber)
}

func EthProfitForBlock(client *ethclient.Client, blocknumber uint, walletAddressString string) *big.Int {
	blocknumberBefore := int64(blocknumber - 1)
	blocknumberAfter := int64(blocknumber)

	balanceBefore, err := client.BalanceAt(context.Background(), common.HexToAddress(walletAddressString), big.NewInt(blocknumberBefore))
	if err != nil {
		log.Fatal("Could not load balance for " + walletAddressString + " for block " + fmt.Sprint(blocknumber))
	}
	balanceAfter, err := client.BalanceAt(context.Background(), common.HexToAddress(walletAddressString), big.NewInt(blocknumberAfter))
	if err != nil {
		log.Fatal("Could not load balance for " + walletAddressString + " for block " + fmt.Sprint(blocknumber))
	}
	profitPerBlock := new(big.Int).Set(balanceAfter).Sub(balanceAfter, balanceBefore)
	txFees := getTransactionFeesByBlock(client, blocknumber, walletAddressString)

	return profitPerBlock.Sub(profitPerBlock, txFees)

}

// TX Fees: Gas Used * Gas Price + Value
func getTransactionFeesByBlock(client *ethclient.Client, blocknumber uint, walletAddressString string) *big.Int {
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

func ConvertWEIToETH(wei *big.Int, precision uint) string {
	weiFloat := new(big.Float).SetInt(wei)
	ethFloat := new(big.Float).Quo(weiFloat, big.NewFloat(1e18))
	eth, _ := ethFloat.Float64()
	return fmt.Sprintf("%.*f", precision, eth)
}
