
## `github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/agreements` Documentation

The `agreements` SDK allows for interaction with Azure Resource Manager `datadog` (API Version `2021-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/agreements"
```


### Client Initialization

```go
client := agreements.NewAgreementsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AgreementsClient.MarketplaceAgreementsCreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := agreements.DatadogAgreementResource{
	// ...
}


read, err := client.MarketplaceAgreementsCreateOrUpdate(ctx, id, payload)
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
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.MarketplaceAgreementsList(ctx, id)` can be used to do batched pagination
items, err := client.MarketplaceAgreementsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
