package dashboard

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/shadmau/MEVProfitExplorer/pkg/ethereum"
)

const DEFAULT_BLOCKS = 500

func ShowTextDashboard(startBlock uint, endBlock uint, rpc string, ignoreWithdrawls bool) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
		return
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

}
