
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/dbnodes` Documentation

The `dbnodes` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2024-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/dbnodes"
```


### Client Initialization

```go
client := dbnodes.NewDbNodesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DbNodesClient.Action`

```go
ctx := context.TODO()
id := dbnodes.NewDbNodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudVmClusterName", "dbNodeName")

payload := dbnodes.DbNodeAction{
	// ...
}


if err := client.ActionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DbNodesClient.Get`

```go
ctx := context.TODO()
id := dbnodes.NewDbNodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudVmClusterName", "dbNodeName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DbNodesClient.ListByCloudVMCluster`

```go
ctx := context.TODO()
id := dbnodes.NewCloudVMClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudVmClusterName")

// alternatively `client.ListByCloudVMCluster(ctx, id)` can be used to do batched pagination
items, err := client.ListByCloudVMClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
