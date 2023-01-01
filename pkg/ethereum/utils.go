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

func GetCurrentBlockNumber(client *ethclient.Client) uint {
	blocknumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal("Could not get Blocknumber!")
	}
	return uint(blocknumber)
}

func EthProfitForBlock(client *ethclient.Client, blocknumber uint, walletAddressString string) big.Int {
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
	result := new(big.Int).Set(balanceAfter).Sub(balanceAfter, balanceBefore)

	//reduce result by TransactionFees
	return *result

}

// todo implement
func getTransactionFeesByBlock(client *ethclient.Client, blocknumber uint, walletAddressString string) *big.Int {

	return big.NewInt(0)
}

func ConvertWEIToETH(wei *big.Int) string {
	weiFloat := new(big.Float).SetInt(wei)
	ethFloat := new(big.Float).Quo(weiFloat, big.NewFloat(1e18))
	eth, _ := ethFloat.Float64()
	return fmt.Sprintf("%.5f", eth)
}
