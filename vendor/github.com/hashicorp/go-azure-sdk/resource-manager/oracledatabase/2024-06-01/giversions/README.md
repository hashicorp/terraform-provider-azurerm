
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/giversions` Documentation

The `giversions` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2024-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/giversions"
```


### Client Initialization

```go
client := giversions.NewGiVersionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GiVersionsClient.Get`

```go
ctx := context.TODO()
id := giversions.NewGiVersionID("12345678-1234-9876-4563-123456789012", "locationName", "giVersionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GiVersionsClient.ListByLocation`

```go
ctx := context.TODO()
id := giversions.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.ListByLocation(ctx, id)` can be used to do batched pagination
items, err := client.ListByLocationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
