
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/locationbasedcapabilityset` Documentation

The `locationbasedcapabilityset` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2024-12-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/locationbasedcapabilityset"
```


### Client Initialization

```go
client := locationbasedcapabilityset.NewLocationBasedCapabilitySetClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LocationBasedCapabilitySetClient.Get`

```go
ctx := context.TODO()
id := locationbasedcapabilityset.NewCapabilitySetID("12345678-1234-9876-4563-123456789012", "locationName", "capabilitySetName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocationBasedCapabilitySetClient.List`

```go
ctx := context.TODO()
id := locationbasedcapabilityset.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
