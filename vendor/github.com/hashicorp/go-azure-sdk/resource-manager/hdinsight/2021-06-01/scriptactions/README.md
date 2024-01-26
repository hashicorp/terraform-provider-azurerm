
## `github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/scriptactions` Documentation

The `scriptactions` SDK allows for interaction with the Azure Resource Manager Service `hdinsight` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/scriptactions"
```


### Client Initialization

```go
client := scriptactions.NewScriptActionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ScriptActionsClient.Delete`

```go
ctx := context.TODO()
id := scriptactions.NewScriptActionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "scriptActionValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScriptActionsClient.ListByCluster`

```go
ctx := context.TODO()
id := scriptactions.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

// alternatively `client.ListByCluster(ctx, id)` can be used to do batched pagination
items, err := client.ListByClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
