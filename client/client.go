package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"path"
	"time"
	"crypto/x509"
)

const (
	channelName   = "assetchannel"
	chaincodeName = "assetcc"
	mspID         = "Org1MSP"
	cryptoPath    = "../../Documents/fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
	certPath      = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem"
	keyPath       = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
	tlsCertPath   = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint  = "localhost:7051"
	gatewayPeer   = "peer0.org1.example.com"
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

// TODO: Create a new investor

func CreateInvestor(contract *client.Contract, investor Investor) error {

	investorJSON, err := json.Marshal(investor)
	if err != nil {
		return fmt.Errorf("failed to marshal the investor details: %v", err)
	}

	_, err = contract.SubmitTransaction("CreateUser", string(investorJSON))
	if err != nil {
		return fmt.Errorf("failed to create the investor: %v", err)
	}

	fmt.Printf("successfully created the investor")
	return nil

}

// TODO: Register a new asset

func RegisterAsset(contract *client.Contract, asset Asset) error {

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal asset details: %v", err)
	}

	_, err = contract.SubmitTransaction("RegisterAsset", string(assetJSON))
	if err != nil {
		return fmt.Errorf("failed to register the asset details: %v", err)
	}

	fmt.Printf("successfully registered the asset details")
	return nil
}

// TODO: Subscribe to the asset

func SubscribeAsset(contract *client.Contract, investorID string, isin string, units int) error {

	_, err := contract.SubmitTransaction("SubscribeAsset", investorID, isin, fmt.Sprint("%d", units))
	if err != nil {
		return fmt.Errorf("failed to subscribe to the asset: %v", err)
	}

	fmt.Printf("successfully subscribed to the asset")
	return nil
}

// TODO: Redeem assets

func RedeemAsset(contract *client.Contract, investorID string, isin string, units int) error {

	_, err := contract.SubmitTransaction("RedeemAsset", investorID, isin, fmt.Sprint("%d", units))
	if err != nil {
		return fmt.Errorf("failed to redeem asset details: %v", err)
	}

	fmt.Printf("successfully redeemed the asset: %s\n")
	return nil
}

// TODO: Check the investor's portfolio balance

func GetPortfolio(contract *client.Contract, investorID string) error {

	result, err := contract.EvaluateTransaction("GetPortfolio", investorID)
	if err != nil {
		return fmt.Errorf("failed to get portfolio details: %v", err)
	}

	fmt.Printf("successfully retrieved the portfolio details %s\n", string(result))
	return nil
}

// Listen to chaincode events

func listenToEvents(eventChannel <-chan *client.ChaincodeEvent) {
	// Register a subscription to chaincode events
	fmt.Println("Listening to chaincode events...")

	for {
		select {
		case ccEvent := <-eventChannel:
			fmt.Printf("Received chaincode event: %s\n", ccEvent.Payload)
		case <-time.After(10 * time.Second):
			fmt.Println("Timeout: No events received for 10 seconds")
			return
		}
	}
}

func main() {
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gateway.Close()

	network := gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)
	eventChannel, err := network.ChaincodeEvents(context.Background(), chaincodeName)
	if err != nil {
		panic(err)
	}

	go listenToEvents(eventChannel)

	// Create a new investor
	investor := Investor{
		InvestorID: "investor1",
		Balance:    10000,
	}
	err = CreateInvestor(contract, investor)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Register a new asset
	asset := Asset{
		ISIN:           "US1234567890",
		CompanyName:    "Example Corp",
		AssetType:      "Stock",
		TotalUnits:     1000,
		PricePerUnit:   100,
		AvailableUnits: 1000,
	}
	err = RegisterAsset(contract, asset)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Subscribe to an asset
	err = SubscribeAsset(contract, "investor1", "US1234567890", 10)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Redeem asset units
	err = RedeemAsset(contract, "investor1", "US1234567890", 5)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get investor's portfolio
	err = GetPortfolio(contract, "investor1")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}
