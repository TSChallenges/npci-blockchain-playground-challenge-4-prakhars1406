# NPCI-Blockchain-Playground-Challenge-4 : Digital Stock Exchange on Hyperledger Fabric

## Description
The objective of this project is to build a blockchain-based broker system using Hyperledger Fabric. The system will facilitate digital asset creation, subscriptions, redemptions, and investor balances between a Broker (i.e.Groww) and an Exchange (NSE or BSE).

## Getting Started

**1. Install Dependencies**:

Ensure the following dependencies are installed:
- Docker and Docker Compose
- Go (latest version)
- Hyperledger Fabric CLI tools
  
**2. Start the Network**:

Run the following command to set up the Fabric network and create a **channel `(assetchannel)`**:
```bash
./network.sh up createChannel -c assetchannel -ca
```

**3. Deploy Chaincode**:

Deploy the **Asset Management** smart contract **`(assetcc)`**:
```bash
./network.sh deployCC -ccn assetcc -ccp ./chaincode -ccl go
```

**4. Run the Client Application**

Start the client application to interact with the blockchain network:
```bash
go run client/client.go
```

## Functionality

### Asset Lifecycle :

- **Create an Asset:** Register a new asset with properties such as **ISIN, company name, asset type, total units, and price per unit**.
- **Subscribe to Asset Units:** Investors can purchase asset units, subject to validation checks.
- **Redeem Asset Units:** Investors can redeem units if they hold sufficient balance.
- **Query Asset Details:** Retrieve asset information from the ledger.
  
### User Roles :

- **Investor:** Can subscribe to and redeem asset units.
- **Asset Manager:** Manages asset creation and lifecycle.
- **Admin:** Manages system configurations and approvals.
  
### Event Emission :

- **Subscription Event:** Triggered when an investor successfully subscribes to asset units.
- **Redemption Event:** Triggered when an investor redeems asset units.


## Validations

### Subscription Validations :

- **Asset Existence Validation:** Ensure the asset exists before allowing subscription.
- **Maximum Subscription Limit:** Ensure the requested units do not exceed the maximum allowed.
- **Available Units Validation:** Ensure enough asset units are available for subscription.
- **Investor Balance Validation:** Verify that the investor has sufficient funds for the subscription.
- **Role-Based Access Control:** Ensure only Investors can subscribe.
  
### Redemption Validations :

- **Asset Existence Validation:** Ensure the asset exists before allowing redemption.
- **Investor Holdings Validation:** Ensure the investor holds enough asset units to redeem.
- **Minimum Redemption Limit:** Validate that the requested units meet the minimum redemption threshold. i.e. 30 units
- **Lock-In Period Validation:** Prevent redemption if the lock-in period has not ended. i.e. 7 days
- **Role-Based Access Control:** Ensure only Investors can redeem.

## TODO
Refer to `TODO.md` for remaining implementation tasks

## Submission
Commit your code to the GitHub repository and create a Pull Request (PR) for review.
