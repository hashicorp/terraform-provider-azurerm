
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/usages` Documentation

The `usages` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/usages"
```


### Client Initialization

```go
client := usages.NewUsagesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `UsagesClient.ListByLocation`

```go
ctx := context.TODO()
id := usages.NewLocationID("12345678-1234-9876-4563-123456789012", "location")

// alternatively `client.ListByLocation(ctx, id)` can be used to do batched pagination
items, err := client.ListByLocationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
