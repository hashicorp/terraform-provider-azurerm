
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2021-05-01/nodetype` Documentation

The `nodetype` SDK allows for interaction with the Azure Resource Manager Service `servicefabricmanagedcluster` (API Version `2021-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2021-05-01/nodetype"
```


### Client Initialization

```go
client := nodetype.NewNodeTypeClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
if err != nil {
	// handle the error
}
```


### Example Usage: `NodeTypeClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := nodetype.NewNodeTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "nodeTypeValue")

payload := nodetype.NodeType{
	// ...
}

future, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `NodeTypeClient.Delete`

```go
ctx := context.TODO()
id := nodetype.NewNodeTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "nodeTypeValue")
future, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `NodeTypeClient.DeleteNode`

```go
ctx := context.TODO()
id := nodetype.NewNodeTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "nodeTypeValue")

payload := nodetype.NodeTypeActionParameters{
	// ...
}

future, err := client.DeleteNode(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `NodeTypeClient.Get`

```go
ctx := context.TODO()
id := nodetype.NewNodeTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "nodeTypeValue")
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
id := nodetype.NewManagedClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")
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
id := nodetype.NewNodeTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "nodeTypeValue")

payload := nodetype.NodeTypeActionParameters{
	// ...
}

future, err := client.Reimage(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `NodeTypeClient.Restart`

```go
ctx := context.TODO()
id := nodetype.NewNodeTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "nodeTypeValue")

payload := nodetype.NodeTypeActionParameters{
	// ...
}

future, err := client.Restart(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `NodeTypeClient.Update`

```go
ctx := context.TODO()
id := nodetype.NewNodeTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "nodeTypeValue")

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
