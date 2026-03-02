
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2024-04-01/nodetype` Documentation

The `nodetype` SDK allows for interaction with Azure Resource Manager `servicefabricmanagedcluster` (API Version `2024-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2024-04-01/nodetype"
```


### Client Initialization

```go
client := nodetype.NewNodeTypeClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NodeTypeClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := nodetype.NewNodeTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName", "nodeTypeName")

payload := nodetype.NodeType{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NodeTypeClient.Delete`

```go
ctx := context.TODO()
id := nodetype.NewNodeTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName", "nodeTypeName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NodeTypeClient.DeleteNode`

```go
ctx := context.TODO()
id := nodetype.NewNodeTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName", "nodeTypeName")

payload := nodetype.NodeTypeActionParameters{
	// ...
}


if err := client.DeleteNodeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NodeTypeClient.Get`

```go
ctx := context.TODO()
id := nodetype.NewNodeTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName", "nodeTypeName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NodeTypeClient.ListByManagedClusters`

```go
ctx := context.TODO()
id := nodetype.NewManagedClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName")

// alternatively `client.ListByManagedClusters(ctx, id)` can be used to do batched pagination
items, err := client.ListByManagedClustersComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NodeTypeClient.Reimage`

```go
ctx := context.TODO()
id := nodetype.NewNodeTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName", "nodeTypeName")

payload := nodetype.NodeTypeActionParameters{
	// ...
}


if err := client.ReimageThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NodeTypeClient.Restart`

```go
ctx := context.TODO()
id := nodetype.NewNodeTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName", "nodeTypeName")

payload := nodetype.NodeTypeActionParameters{
	// ...
}


if err := client.RestartThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NodeTypeClient.SkusList`

```go
ctx := context.TODO()
id := nodetype.NewNodeTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName", "nodeTypeName")

// alternatively `client.SkusList(ctx, id)` can be used to do batched pagination
items, err := client.SkusListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NodeTypeClient.Update`

```go
ctx := context.TODO()
id := nodetype.NewNodeTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName", "nodeTypeName")

payload := nodetype.NodeTypeUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
