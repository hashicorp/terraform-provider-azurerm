
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/reachabilityanalysisruns` Documentation

The `reachabilityanalysisruns` SDK allows for interaction with Azure Resource Manager `network` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/reachabilityanalysisruns"
```


### Client Initialization

```go
client := reachabilityanalysisruns.NewReachabilityAnalysisRunsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReachabilityAnalysisRunsClient.Create`

```go
ctx := context.TODO()
id := reachabilityanalysisruns.NewReachabilityAnalysisRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "verifierWorkspaceName", "reachabilityAnalysisRunName")

payload := reachabilityanalysisruns.ReachabilityAnalysisRun{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReachabilityAnalysisRunsClient.Delete`

```go
ctx := context.TODO()
id := reachabilityanalysisruns.NewReachabilityAnalysisRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "verifierWorkspaceName", "reachabilityAnalysisRunName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReachabilityAnalysisRunsClient.Get`

```go
ctx := context.TODO()
id := reachabilityanalysisruns.NewReachabilityAnalysisRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "verifierWorkspaceName", "reachabilityAnalysisRunName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReachabilityAnalysisRunsClient.List`

```go
ctx := context.TODO()
id := reachabilityanalysisruns.NewVerifierWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "verifierWorkspaceName")

// alternatively `client.List(ctx, id, reachabilityanalysisruns.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, reachabilityanalysisruns.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
