
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/logicalnetworks` Documentation

The `logicalnetworks` SDK allows for interaction with Azure Resource Manager `azurestackhci` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/logicalnetworks"
```


### Client Initialization

```go
client := logicalnetworks.NewLogicalNetworksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LogicalNetworksClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := logicalnetworks.NewLogicalNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "logicalNetworkName")

payload := logicalnetworks.LogicalNetworks{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LogicalNetworksClient.Delete`

```go
ctx := context.TODO()
id := logicalnetworks.NewLogicalNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "logicalNetworkName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LogicalNetworksClient.Get`

```go
ctx := context.TODO()
id := logicalnetworks.NewLogicalNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "logicalNetworkName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LogicalNetworksClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LogicalNetworksClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LogicalNetworksClient.Update`

```go
ctx := context.TODO()
id := logicalnetworks.NewLogicalNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "logicalNetworkName")

payload := logicalnetworks.LogicalNetworksUpdateRequest{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
