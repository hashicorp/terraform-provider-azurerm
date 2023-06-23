
## `github.com/hashicorp/go-azure-sdk/resource-manager/marketplaceordering/2015-06-01/agreements` Documentation

The `agreements` SDK allows for interaction with the Azure Resource Manager Service `marketplaceordering` (API Version `2015-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/marketplaceordering/2015-06-01/agreements"
```


### Client Initialization

```go
client := agreements.NewAgreementsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AgreementsClient.MarketplaceAgreementsCancel`

```go
ctx := context.TODO()
id := agreements.NewPlanID("12345678-1234-9876-4563-123456789012", "publisherIdValue", "offerIdValue", "planIdValue")

read, err := client.MarketplaceAgreementsCancel(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AgreementsClient.MarketplaceAgreementsCreate`

```go
ctx := context.TODO()
id := agreements.NewOfferPlanID("12345678-1234-9876-4563-123456789012", "publisherIdValue", "offerIdValue", "planIdValue")

payload := agreements.AgreementTerms{
	// ...
}


read, err := client.MarketplaceAgreementsCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AgreementsClient.MarketplaceAgreementsGet`

```go
ctx := context.TODO()
id := agreements.NewOfferPlanID("12345678-1234-9876-4563-123456789012", "publisherIdValue", "offerIdValue", "planIdValue")

read, err := client.MarketplaceAgreementsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AgreementsClient.MarketplaceAgreementsGetAgreement`

```go
ctx := context.TODO()
id := agreements.NewPlanID("12345678-1234-9876-4563-123456789012", "publisherIdValue", "offerIdValue", "planIdValue")

read, err := client.MarketplaceAgreementsGetAgreement(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AgreementsClient.MarketplaceAgreementsList`

```go
ctx := context.TODO()
id := agreements.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.MarketplaceAgreementsList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AgreementsClient.MarketplaceAgreementsSign`

```go
ctx := context.TODO()
id := agreements.NewPlanID("12345678-1234-9876-4563-123456789012", "publisherIdValue", "offerIdValue", "planIdValue")

read, err := client.MarketplaceAgreementsSign(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
