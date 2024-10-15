
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/locationbasedcapabilities` Documentation

The `locationbasedcapabilities` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2022-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/locationbasedcapabilities"
```


### Client Initialization

```go
client := locationbasedcapabilities.NewLocationBasedCapabilitiesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LocationBasedCapabilitiesClient.List`

```go
ctx := context.TODO()
id := locationbasedcapabilities.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
