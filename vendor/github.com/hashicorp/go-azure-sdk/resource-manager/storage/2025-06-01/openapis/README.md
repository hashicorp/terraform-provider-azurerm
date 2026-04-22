
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/openapis` Documentation

The `openapis` SDK allows for interaction with Azure Resource Manager `storage` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/openapis"
```


### Client Initialization

```go
client := openapis.NewOpenapisClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OpenapisClient.DeletedAccountsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.DeletedAccountsList(ctx, id)` can be used to do batched pagination
items, err := client.DeletedAccountsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.SkusList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.SkusList(ctx, id)` can be used to do batched pagination
items, err := client.SkusListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.UsagesListByLocation`

```go
ctx := context.TODO()
id := openapis.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.UsagesListByLocation(ctx, id)` can be used to do batched pagination
items, err := client.UsagesListByLocationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
