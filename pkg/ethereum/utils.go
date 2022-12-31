package ethereum

import (
	"context"
	"log"

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
