package services

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ruhil6789/event-sky/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MyService struct {
	Client *ethclient.Client
}

const (
	url             = "https://bsc-testnet.publicnode.com"
	contractAddress = "0x243C43588fb01AC9667eeA8e041F90E4326FBc15"
	// abiPath="./contracts/ERC20.json"
)

// StakedEvent represents the structure of the Staked event from the smart contract.
type StakedEvent struct {
	User            common.Address // address user (index_topic_1)
	Referrer        common.Address // address referrer (index_topic_2)
	DepositTime     *big.Int       // uint256 depositTime
	Amount          *big.Int       // uint256 amount
	TotalInvestment *big.Int       // uint256 totalInvestment
}

func (c *MyService) getEvents(address string, chain string, contractABI abi.ABI) ([]types.Log, error) {

	client, err := ethclient.Dial("https://bsc-testnet.publicnode.com")
	log.Println(client, "client")
	if err != nil {
		log.Fatal((err))
	}

	fmt.Println("connected to Node RPC ")

	var eventBatchSize int64

	if chain == "Bsc" {
		eventBatchSize = big.NewInt(1000).Int64()
	} else {
		eventBatchSize = big.NewInt(100).Int64()
	}

	currentBlock, err := c.Client.BlockNumber(context.Background())

	fmt.Println(currentBlock, "current block no")
	if err != nil {
		log.Printf("Failed to get block number %v", err)
		return nil, err
	}

	startBlock, err := getBlockInfo(address, chain)

	if err != nil || startBlock == 0 {
		log.Printf("Failed to get start block: %v", err)
		return nil, err
	}

	var endBlock uint64
	if startBlock+uint64(eventBatchSize) > currentBlock {
		endBlock = currentBlock
	} else {
		endBlock = startBlock + uint64(eventBatchSize)
	}

	if currentBlock > startBlock {
		contractAddress := common.HexToAddress(strings.ToLower(address))

		// contractInstance := client.FilterLogs()
		contractABI := `[
	 {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "internalType": "address",
                "name": "user",
                "type": "address"
            },
            {
                "indexed": true,
                "internalType": "address",
                "name": "referrer",
                "type": "address"
            },
            {
                "indexed": false,
                "internalType": "uint256",
                "name": "depositTime",
                "type": "uint256"
            },
            {
                "indexed": false,
                "internalType": "uint256",
                "name": "amount",
                "type": "uint256"
            },
            {
                "indexed": false,
                "internalType": "uint256",
                "name": "totalInvestment",
                "type": "uint256"
            }
        ],
        "name": "Staked",
        "type": "event"
     }
	 ]`

		parsedAbi, err := abi.JSON(strings.NewReader(contractABI))
		if err != nil {
			log.Fatal(err)
		}
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(startBlock)),
			ToBlock:   big.NewInt(int64(endBlock)),
			Addresses: []common.Address{contractAddress},
			Topics: [][]common.Hash{
				// {common.HexToHash("0x6997c1b5b980012a5c54a08f8852865180742947186756624460059058236934")}, // transfer
				{
					parsedAbi.Events["Staked"].ID,
					parsedAbi.Events["thresholdUpdated"].ID,
				},
			},
		}

		logs, err := c.Client.FilterLogs(context.Background(), query)
		if err != nil {
			log.Printf("Failed to fetch events: %v", err)
			return nil, err
		}

		for _, vLog := range logs {
			switch vLog.Topics[0].Hex() {
			case parsedAbi.Events["Staked"].ID.Hex():
				event := StakedEvent{}
				err := parsedAbi.UnpackIntoInterface(&event, "Staked", vLog.Data)
				if err != nil {
					log.Fatal(err)
				}
				event.User = common.HexToAddress(vLog.Topics[1].Hex())
				event.Referrer = common.HexToAddress(vLog.Topics[2].Hex())
				depositTime := common.HexToAddress(vLog.Topics[3].Hex())
				event.DepositTime = big.NewInt(0).SetBytes(depositTime.Bytes())
				// Extract numeric fields from data
				// event.DepositTime = new(big.Int).SetBytes(vLog.Topics[4])       // Assuming first 32 bytes
				event.Amount = new(big.Int).SetBytes(vLog.Data[32:64])          // Next 32 bytes
				event.TotalInvestment = new(big.Int).SetBytes(vLog.Data[64:96]) // Next 32 bytes

				// Print out the event details
				fmt.Printf("Staked: User: %s, Referrer: %s, Deposit Time: %s, Amount: %s, Total Investment: %s\n",
					event.User.Hex(),
					event.Referrer.Hex(),
					event.DepositTime.String(),
					event.Amount.String(),
					event.TotalInvestment.String())

				fmt.Printf("Staked: %s tokens from %s to %s\n", event.TotalInvestment.String(),
					event.Amount.String(), event.User.Hex(), event.Referrer.Hex())

			default:
				fmt.Println("Unknown event:", vLog.Topics[0].Hex())
			}

			//      events :=getPastEvents(address,chain,endBlock)

			if len(logs) > 0 {
				updateErr := updateBlockInfo(address, chain, endBlock)
				if updateErr != nil {
					log.Printf("Failed to update block info: %v", updateErr)
					return nil, updateErr
				}
			}
			return logs, nil

		}
	}
	log.Println("Waiting for block number to update")
	return nil, nil

}

var client *mongo.Client

func (c *MyService) getPastEvents(contractAddress common.Address, fromBlock, toBlock *big.Int) ([]types.Log, error) {
	// Define the query parameters
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []common.Address{contractAddress},
	}

	// Fetch the past events
	logs, err := c.Client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Printf("Error fetching past events: %v", err)
		return nil, err
	}

	return logs, nil
}

func updateBlockInfo(address, chain string, endBlock uint64) error {
	collection := client.Database("event-sky").Collection("eventsModel")

	filter := bson.M{
		"address": address,
		"chain":   chain,
	}
	update := bson.M{
		"$set": bson.M{"blockNumber": endBlock + 1},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)

	if err != nil {
		log.Printf("Failed to update block info: %v", err)
		return err
	}
	return nil
}

// func (c *MyService) updateBlockInfo(address string, chain string, blockNumber int64) error {

// 	collection := client.Database("event-sky").Collection("eventsModel")

// 	filter := bson.M{
// 		"address": address,
// 		"chain":   chain,
// 	}
// 	update := bson.M{
// 		"$set": bson.M{"blockNumber": blockNumber + 1},
// 	}

// 	opts := options.Update().SetUpsert(true)
// 	_, err := collection.UpdateOne(context.Background(), filter, update, opts)

// 	if err != nil {
// 		log.Printf("Failed to update block info: %v", err)
// 		return err
// 	}
//    return nil
// }

func getBlockInfo(address string, chain string) (uint64, error) {

	collection := client.Database("event-sky").Collection("eventsModel")

	var blockInfo model.EventBlock
	filter := bson.M{
		"address": address,
		"chain":   chain,
	}
	var startBlock uint64 = 41958006
	err := collection.FindOne(context.TODO(), filter).Decode(&blockInfo)
	if err != nil {
		if err == mongo.ErrNilDocument {
			return uint64(startBlock), nil
		}
	}

	return uint64(*blockInfo.BlockNumber), nil

}

func web3Service() {

	client, err := ethclient.Dial("https://bsc-testnet.publicnode.com")
	log.Println(client, "client")
	if err != nil {
		log.Fatal((err))
	}

	fmt.Println("connected to Node RPC ")

}
