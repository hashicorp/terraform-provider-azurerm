
## `github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/nodepools` Documentation

The `nodepools` SDK allows for interaction with Azure Resource Manager `discovery` (API Version `2026-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/nodepools"
```


### Client Initialization

```go
client := nodepools.NewNodePoolsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NodePoolsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := nodepools.NewNodePoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "supercomputerName", "nodePoolName")

payload := nodepools.NodePool{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NodePoolsClient.Delete`

```go
ctx := context.TODO()
id := nodepools.NewNodePoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "supercomputerName", "nodePoolName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NodePoolsClient.Get`

```go
ctx := context.TODO()
id := nodepools.NewNodePoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "supercomputerName", "nodePoolName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NodePoolsClient.ListBySupercomputer`

```go
ctx := context.TODO()
id := nodepools.NewSupercomputerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "supercomputerName")

// alternatively `client.ListBySupercomputer(ctx, id)` can be used to do batched pagination
items, err := client.ListBySupercomputerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NodePoolsClient.Update`

```go
ctx := context.TODO()
id := nodepools.NewNodePoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "supercomputerName", "nodePoolName")

payload := nodepools.NodePoolUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
