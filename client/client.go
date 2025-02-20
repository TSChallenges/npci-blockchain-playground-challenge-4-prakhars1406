package main

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"os"
)

func main() {
	// Set up fabric connection
	wallet, err := gateway.NewFileSystemWallet("./wallet")
	if err != nil {
		fmt.Println("Failed to create wallet: ", err)
		os.Exit(1)
	}

	// Get gateway connection
	gw, err := gateway.Connect(
		gateway.WithAddress("localhost:7051"),
		gateway.WithIdentity(wallet, "admin"),
	)
	if err != nil {
		fmt.Println("Failed to connect to gateway: ", err)
		os.Exit(1)
	}

	// Create a client for interacting with the blockchain
	channelClient, err := channel.New(gw, channel.WithChannelID("assetchannel"))
	if err != nil {
		fmt.Println("Failed to create channel client: ", err)
		os.Exit(1)
	}

	// TODO: Create a new investor
	
	
	// TODO: Register a new asset
	
	
	// TODO: Subscribe to the asset
	
	
	// TODO: Redeem assets
	
	
	// TODO: Check the investor's portfolio balance


}
