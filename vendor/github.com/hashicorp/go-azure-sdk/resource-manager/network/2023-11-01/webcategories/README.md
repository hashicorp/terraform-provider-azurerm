
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/webcategories` Documentation

The `webcategories` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/webcategories"
```


### Client Initialization

```go
client := webcategories.NewWebCategoriesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `WebCategoriesClient.Get`

```go
ctx := context.TODO()
id := webcategories.NewAzureWebCategoryID("12345678-1234-9876-4563-123456789012", "azureWebCategoryValue")

read, err := client.Get(ctx, id, webcategories.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebCategoriesClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
