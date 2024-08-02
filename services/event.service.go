package services

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ruhil6789/event-sky/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func filterStakedEvents(data []types.Log, contractABI abi.ABI) []types.Log {
	stakedEvents := make([]types.Log, 0)
	stakedEventID := contractABI.Events["Staked"].ID

	for _, element := range data {
		if element.Topics[0] == stakedEventID {
			stakedEvents = append(stakedEvents, element)
		}
	}

	return stakedEvents
}

func fetch(address string, chain string, abi abi.ABI) {

	data := services.getPastEvents()
	log.Println(data, "data")

	if len(data) > 0 {
    log.Println("Fetched data:", data)

    stakedEvents := filterStakedEvents(data, abi)
    if len(stakedEvents) > 0 {
        err := stakedFunction(stakedEvents, chain, abi)
        if err != nil {
            log.Printf("Error processing staked events: %v", err)
            // Handle the error as needed (e.g., return from the function)
        }
    } else {
        log.Println("No staked events found")
    }
} else {
    log.Println("No data fetched")
}
	
}

 var  db *mongo.Database

func stakedFunction(events []types.Log, chain string, contractABI abi.ABI) error {
    ctx := context.Background()

    // Process events
    var txHistoryData []interface{}
    var usersData []mongo.WriteModel
    var referrerUsersData []mongo.WriteModel

    for _, event := range events {
        var stakedEvent StakedEvent
        err := contractABI.UnpackIntoInterface(&stakedEvent, "Staked", event.Data)
        if err != nil {
            return fmt.Errorf("failed to unpack event: %v", err)
        }

        // Prepare transaction history data
        txHistoryData = append(txHistoryData, bson.M{
            "receiverAddress": strings.ToLower(stakedEvent.User.Hex()),
            "amount":          stakedEvent.Amount.String(),
            "amt":             new(big.Float).Quo(new(big.Float).SetInt(stakedEvent.Amount), new(big.Float).SetInt(big.NewInt(1e18))).String(),
            "event":           "Staked",
            "timestamp":       stakedEvent.DepositTime.Uint64(),
            "transactionHash": event.TxHash.Hex(),
            "blockNumber":     event.BlockNumber,
        })

        // Prepare users data
        newObjectID := primitive.NewObjectID()
        usersData = append(usersData, mongo.NewUpdateOneModel().
            SetFilter(bson.M{"walletAddress": strings.ToLower(stakedEvent.User.Hex())}).
            SetUpdate(bson.M{
                "$set": bson.M{
                    "timestamp":    stakedEvent.DepositTime.Uint64(),
                    "totalStaked":  stakedEvent.TotalInvestment.String(),
                    "packageAmount": stakedEvent.Amount.String(),
                },
                "$setOnInsert": bson.M{
                    "_id":             newObjectID,
                    "referrerAddress": strings.ToLower(stakedEvent.Referrer.Hex()),
                    "referralLink":    fmt.Sprintf("https://%s/auth/funds/%s","abcdef", newObjectID.Hex()),
                    "countOfReferree": 0,
                    "claimableAmount": 0,
                },
            }).
            SetUpsert(true))

        if strings.ToLower(chain) == "bsc" {
            referrerUsersData = append(referrerUsersData, mongo.NewUpdateOneModel().
                SetFilter(bson.M{"walletAddress": strings.ToLower(stakedEvent.Referrer.Hex())}).
                SetUpdate(bson.M{
                    "$inc": bson.M{"countOfReferree": 1},
                }))
        }
    }

    // Insert transaction history
    _, err := db.Collection("allTx").InsertMany(ctx, txHistoryData, options.InsertMany().SetOrdered(false))
    if err != nil {
        return fmt.Errorf("failed to insert transaction history: %v", err)
    }

    // Update users
    _, err = db.Collection("users").BulkWrite(ctx, usersData, options.BulkWrite().SetOrdered(false))
    if err != nil {
        return fmt.Errorf("failed to update users: %v", err)
    }

    // Update referrer users (BSC only)
    if strings.ToLower(chain) == "bsc" && len(referrerUsersData) > 0 {
        _, err = db.Collection("users").BulkWrite(ctx, referrerUsersData, options.BulkWrite().SetOrdered(false))
        if err != nil {
            return fmt.Errorf("failed to update referrer users: %v", err)
        }
    }

    return nil
}