package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
)

type Asset struct {
	ISIN           string `json:"isin"`
	CompanyName    string `json:"company_name"`
	AssetType      string `json:"asset_type"`
	TotalUnits     int    `json:"total_units"`
	PricePerUnit   int    `json:"price_per_unit"`
	AvailableUnits int    `json:"available_units"`
}

type Investor struct {
	InvestorID string `json:"investor_id"`
	Balance    int    `json:"balance"`
}

type AssetManagementContract struct {
	contractapi.Contract
}

// TODO: CreateUser function adds an investor to the ledger


// TODO: RegisterAsset function registers a new asset


// TODO: SubscribeAsset function allows an investor to subscribe to asset units


// TODO: RedeemAsset function allows an investor to redeem asset units


// TODO: GetPortfolio function retrieves an investor's portfolio



func main() {
	chaincode, err := contractapi.NewChaincode(&AssetManagementContract{})
	if err != nil {
		fmt.Printf("Error creating asset management chaincode: %v", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting asset management chaincode: %v", err)
	}
}

