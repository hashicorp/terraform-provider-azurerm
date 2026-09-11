
## `github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogs` Documentation

The `datadogs` SDK allows for interaction with Azure Resource Manager `datadog` (API Version `2025-06-11`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogs"
```


### Client Initialization

```go
client := datadogs.NewDatadogsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DatadogsClient.CreationSupportedGet`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.CreationSupportedGet(ctx, id, datadogs.DefaultCreationSupportedGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatadogsClient.CreationSupportedList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.CreationSupportedList(ctx, id, datadogs.DefaultCreationSupportedListOperationOptions())` can be used to do batched pagination
items, err := client.CreationSupportedListComplete(ctx, id, datadogs.DefaultCreationSupportedListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DatadogsClient.MarketplaceAgreementsCreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := datadogs.DatadogAgreementResource{
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


### Example Usage: `DatadogsClient.MarketplaceAgreementsList`

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
