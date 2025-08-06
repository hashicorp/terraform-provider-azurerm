
## `github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2024-05-01/vnetpeering` Documentation

The `vnetpeering` SDK allows for interaction with Azure Resource Manager `databricks` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2024-05-01/vnetpeering"
```


### Client Initialization

```go
client := vnetpeering.NewVNetPeeringClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VNetPeeringClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := vnetpeering.NewVirtualNetworkPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "virtualNetworkPeeringName")

payload := vnetpeering.VirtualNetworkPeering{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VNetPeeringClient.Delete`

```go
ctx := context.TODO()
id := vnetpeering.NewVirtualNetworkPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "virtualNetworkPeeringName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VNetPeeringClient.Get`

```go
ctx := context.TODO()
id := vnetpeering.NewVirtualNetworkPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "virtualNetworkPeeringName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VNetPeeringClient.ListByWorkspace`

```go
ctx := context.TODO()
id := vnetpeering.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

// alternatively `client.ListByWorkspace(ctx, id)` can be used to do batched pagination
items, err := client.ListByWorkspaceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
