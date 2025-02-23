package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
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
	InvestorID             string           `json:"investor_id"`
	Balance                int              `json:"balance"`
	Holdings               map[string]int   `json:"holdings"`
	SubscriptionTimestamps map[string]int64 `json:"subscriptionTimestamps"`
}

type AssetManagementContract struct {
	contractapi.Contract
}

var (
	LockInTime int64 = 7 * 24 * 60 * 60 // 7 days in seconds
)

type SubscriptionRedeemRequest struct {
	InvestorID string `json:"investorID"`
	ISIN       string `json:"isin"`
	Units      int    `json:"units"`
	Timestamp  int64  `json:"timestamp"`
}

// TODO: CreateUser function adds an investor to the ledger
func (amc *AssetManagementContract) CreateUser(ctx contractapi.TransactionContextInterface, investor Investor) error {
	if investor.InvestorID == "" {
		fmt.Errorf("investor id cannot be empty")
		return fmt.Errorf("investor id cannot be empty")
	}

	if investor.Balance <= 0 {
		fmt.Errorf("investor balance cannot be less then 0")
		return fmt.Errorf("investor balance cannot be less then 0")
	}

	data, err := ctx.GetStub().GetState(investor.InvestorID)
	if err == nil || data != nil {
		fmt.Errorf("investor already exists")
		return fmt.Errorf("investor already exists")
	}

	investorDataInBytes, err := json.Marshal(investor)
	if err != nil {
		fmt.Errorf("error in marshalling investor data  %s", err.Error())
		return fmt.Errorf("error in marshalling investor data")
	}

	err = ctx.GetStub().PutState(investor.InvestorID, investorDataInBytes)
	if err != nil {
		fmt.Errorf("error in adding investor %s", err.Error())
		return fmt.Errorf("error in adding investor")
	}

	err = ctx.GetStub().SetEvent("CreateUser", investorDataInBytes)
	if err != nil {
		fmt.Errorf("error in setting event %s", err.Error())
		return err
	}

	return nil
}

// TODO: RegisterAsset function registers a new asset

func (amc *AssetManagementContract) RegisterAsset(ctx contractapi.TransactionContextInterface, asset Asset) error {
	if asset.ISIN == "" {
		fmt.Errorf("asset ISIN cannot be empty")
		return fmt.Errorf("asset ISIN cannot be empty")
	}

	if asset.CompanyName == "" {
		fmt.Errorf("asset CompanyName cannot be empty")
		return fmt.Errorf("asset CompanyName cannot be empty")
	}

	if asset.AssetType == "" {
		fmt.Errorf("asset AssetType cannot be empty")
		return fmt.Errorf("asset AssetType cannot be empty")
	}

	if asset.TotalUnits <= 0 {
		fmt.Errorf("asset TotalUnits cannot be less then 0")
		return fmt.Errorf("asset TotalUnits cannot be less then 0")
	}

	if asset.PricePerUnit <= 0 {
		fmt.Errorf("asset PricePerUnit cannot be less then 0")
		return fmt.Errorf("asset PricePerUnit cannot be less then 0")
	}

	if asset.AvailableUnits <= 0 {
		fmt.Errorf("asset AvailableUnits cannot be less then 0")
		return fmt.Errorf("asset AvailableUnits cannot be less then 0")
	}

	data, err := ctx.GetStub().GetState(asset.ISIN)
	if err == nil || data != nil {
		fmt.Errorf("asset already exists")
		return fmt.Errorf("asset already exists")
	}

	assetInBytes, err := json.Marshal(asset)
	if err != nil {
		fmt.Errorf("error in mashalling asset %s", err.Error())
		return err
	}

	err = ctx.GetStub().PutState(asset.ISIN, assetInBytes)
	if err != nil {
		fmt.Errorf("error in adding asset %s", err.Error())
		return err
	}

	err = ctx.GetStub().SetEvent("RegisterAsset", assetInBytes)
	if err != nil {
		fmt.Errorf("error in setting event %s", err.Error())
		return err
	}

	return nil
}

func (s *AssetManagementContract) SubscribeAsset(ctx contractapi.TransactionContextInterface, requestJSON string) error {

	var request SubscriptionRedeemRequest
	err := json.Unmarshal([]byte(requestJSON), &request)
	if err != nil {
		return fmt.Errorf("failed to unmarshal request: %v", err)
	}

	assetJSON, err := ctx.GetStub().GetState(request.ISIN)
	if err != nil {
		return fmt.Errorf("failed to read asset from world state: %v", err)
	}
	if assetJSON == nil {
		return fmt.Errorf("asset with ISIN %s does not exist", request.ISIN)
	}

	investorJSON, err := ctx.GetStub().GetState(request.InvestorID)
	if err != nil {
		return fmt.Errorf("failed to read investor from world state: %v", err)
	}
	if investorJSON == nil {
		return fmt.Errorf("investor with ID %s does not exist", request.InvestorID)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return err
	}
	var investor Investor
	err = json.Unmarshal(investorJSON, &investor)
	if err != nil {
		return err
	}

	if investor.Holdings == nil {
		investor.Holdings = make(map[string]int)
	}
	if investor.SubscriptionTimestamps == nil {
		investor.SubscriptionTimestamps = make(map[string]int64)
	}

	if request.Units > asset.AvailableUnits {
		return fmt.Errorf("requested units exceed available units for asset %s", request.ISIN)
	}

	totalCost := request.Units * asset.PricePerUnit
	if totalCost > investor.Balance {
		return fmt.Errorf("insufficient balance for investor %s", request.InvestorID)
	}

	asset.AvailableUnits -= request.Units

	investor.Balance -= totalCost
	investor.Holdings[request.ISIN] = request.Units
	investor.SubscriptionTimestamps[request.ISIN] = request.Timestamp

	assetJSON, err = json.Marshal(asset)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(request.ISIN, assetJSON)
	if err != nil {
		return err
	}

	investorJSON, err = json.Marshal(investor)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(request.InvestorID, investorJSON)
	if err != nil {
		return err
	}

	eventPayload := fmt.Sprintf("Investor %s subscribed to %d units of asset %s", request.InvestorID, request.Units, request.ISIN)
	return ctx.GetStub().SetEvent("SubscriptionEvent", []byte(eventPayload))
}

func (s *AssetManagementContract) RedeemAsset(ctx contractapi.TransactionContextInterface, requestJSON string) error {

	var request SubscriptionRedeemRequest
	err := json.Unmarshal([]byte(requestJSON), &request)
	if err != nil {
		return fmt.Errorf("failed to unmarshal request: %v", err)
	}

	assetJSON, err := ctx.GetStub().GetState(request.ISIN)
	if err != nil {
		return fmt.Errorf("failed to read asset from world state: %v", err)
	}
	if assetJSON == nil {
		return fmt.Errorf("asset with ISIN %s does not exist", request.ISIN)
	}

	investorJSON, err := ctx.GetStub().GetState(request.InvestorID)
	if err != nil {
		return fmt.Errorf("failed to read investor from world state: %v", err)
	}
	if investorJSON == nil {
		return fmt.Errorf("investor with ID %s does not exist", request.InvestorID)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return err
	}
	var investor Investor
	err = json.Unmarshal(investorJSON, &investor)
	if err != nil {
		return err
	}

	if investor.Holdings[request.ISIN] < request.Units {
		return fmt.Errorf("insufficient holdings for investor %s", request.InvestorID)
	}

	if request.Units < 30 {
		return fmt.Errorf("minimum redemption limit is 30 units")
	}

	lockInEndTime := investor.SubscriptionTimestamps[request.ISIN] + LockInTime
	if request.Timestamp < lockInEndTime {
		return fmt.Errorf("asset %s is still under lock-in period", request.ISIN)
	}

	asset.AvailableUnits += request.Units
	investor.Balance += request.Units * asset.PricePerUnit
	investor.Holdings[request.ISIN] = investor.Holdings[request.ISIN] - request.Units

	assetJSON, err = json.Marshal(asset)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(request.ISIN, assetJSON)
	if err != nil {
		return err
	}

	investorJSON, err = json.Marshal(investor)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(request.InvestorID, investorJSON)
	if err != nil {
		return err
	}

	eventPayload := fmt.Sprintf("Investor %s redeemed %d units of asset %s", request.InvestorID, request.Units, request.ISIN)
	return ctx.GetStub().SetEvent("RedemptionEvent", []byte(eventPayload))
}

func (s *AssetManagementContract) GetPortfolio(ctx contractapi.TransactionContextInterface, investorID string) (map[string]interface{}, error) {
	investorJSON, err := ctx.GetStub().GetState(investorID)
	if err != nil {
		return nil, fmt.Errorf("failed to read investor from world state: %v", err)
	}
	if investorJSON == nil {
		return nil, fmt.Errorf("investor with ID %s does not exist", investorID)
	}

	var investor Investor
	err = json.Unmarshal(investorJSON, &investor)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal investor: %v", err)
	}

	portfolio := make(map[string]int)
	for isin, units := range investor.Holdings {
		portfolio[isin] = units
	}
	result := map[string]interface{}{
		"balance":   investor.Balance,
		"portfolio": portfolio,
	}

	return result, nil
}
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
