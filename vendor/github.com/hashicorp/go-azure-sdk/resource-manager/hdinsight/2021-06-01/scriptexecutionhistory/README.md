
## `github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/scriptexecutionhistory` Documentation

The `scriptexecutionhistory` SDK allows for interaction with Azure Resource Manager `hdinsight` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/scriptexecutionhistory"
```


### Client Initialization

```go
client := scriptexecutionhistory.NewScriptExecutionHistoryClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ScriptExecutionHistoryClient.ListByCluster`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

// alternatively `client.ListByCluster(ctx, id)` can be used to do batched pagination
items, err := client.ListByClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ScriptExecutionHistoryClient.ScriptActionsGetExecutionDetail`

```go
ctx := context.TODO()
id := scriptexecutionhistory.NewScriptExecutionHistoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "scriptExecutionId")

read, err := client.ScriptActionsGetExecutionDetail(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
