package services

import (
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/robfig/cron/v3"
	"github.com/ruhil6789/event-sky/services"
)
 const contractABI = `[
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
func main() {
    c := cron.New()
	contractAddress:="0x243C43588fb01AC9667eeA8e041F90E4326FBc15"
	Address:= common.Address{contractAddress}

    // Add a cron job
    _, err := c.AddFunc("*/5 * * * *", func() {
		services.fetch("0x243C43588fb01AC9667eeA8e041F90E4326FBc15","Bsc",)
        fmt.Println("Running cron job every minute:", time.Now())
    })
    if err != nil {
        log.Fatal(err)
    }

    // Start the cron scheduler
    c.Start()

    // Run the program indefinitely
    select {}
}