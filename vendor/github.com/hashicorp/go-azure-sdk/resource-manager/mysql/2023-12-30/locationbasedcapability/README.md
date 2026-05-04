
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/locationbasedcapability` Documentation

The `locationbasedcapability` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2023-12-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/locationbasedcapability"
```


### Client Initialization

```go
client := locationbasedcapability.NewLocationBasedCapabilityClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LocationBasedCapabilityClient.SetGet`

```go
ctx := context.TODO()
id := locationbasedcapability.NewCapabilitySetID("12345678-1234-9876-4563-123456789012", "locationName", "capabilitySetName")

read, err := client.SetGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocationBasedCapabilityClient.SetList`

```go
ctx := context.TODO()
id := locationbasedcapability.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.SetList(ctx, id)` can be used to do batched pagination
items, err := client.SetListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
