# TODOs 

## Chaincode (Go)

1. Implement `CreateUser` to add a new investor to the ledger.
2. Implement `RegisterAsset` to create a new asset with properties such as ISIN, company name, asset type, and price per unit.
3. Implement `SubscribeAsset` to allow investors to subscribe to asset units while ensuring validation checks.
4. Implement `RedeemAsset` to allow investors to redeem their asset units while checking available holdings and redemption rules.
5. Implement `GetPortfolio` to retrieve an investor's asset holdings and balance from the ledger.
   
## Client (Go)

1. Call `CreateUser` in `client.go` to register a new investor.
2. Call `RegisterAsset` to create an asset and store its details on the ledger.
3. Add functions for `SubscribeAsset` , `RedeemAsset` , and `GetPortfolio` in the client application.
4. Test interactions with the chaincode by creating an asset, subscribing to units, redeeming units, and checking the investor's portfolio balance.

## Testing
1. Test multiple **asset subscriptions and redemptions** to ensure the ledger accurately reflects ownership changes.
2. Test **edge cases**, such as:
   - Subscribing to more asset units than available.
   - Redeeming more asset units than the investor holds.
   - Checking the portfolio of a non-existent investor.
4. Listen to events in client application
