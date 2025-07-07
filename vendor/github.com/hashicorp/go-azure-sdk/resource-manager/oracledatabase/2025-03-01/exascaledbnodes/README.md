
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/exascaledbnodes` Documentation

The `exascaledbnodes` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2025-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/exascaledbnodes"
```


### Client Initialization

```go
client := exascaledbnodes.NewExascaleDbNodesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExascaleDbNodesClient.Action`

```go
ctx := context.TODO()
id := exascaledbnodes.NewExadbVMClusterDbNodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "exadbVmClusterName", "dbNodeName")

payload := exascaledbnodes.DbNodeAction{
	// ...
}


if err := client.ActionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExascaleDbNodesClient.Get`

```go
ctx := context.TODO()
id := exascaledbnodes.NewExadbVMClusterDbNodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "exadbVmClusterName", "dbNodeName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExascaleDbNodesClient.ListByParent`

```go
ctx := context.TODO()
id := exascaledbnodes.NewExadbVMClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "exadbVmClusterName")

// alternatively `client.ListByParent(ctx, id)` can be used to do batched pagination
items, err := client.ListByParentComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
