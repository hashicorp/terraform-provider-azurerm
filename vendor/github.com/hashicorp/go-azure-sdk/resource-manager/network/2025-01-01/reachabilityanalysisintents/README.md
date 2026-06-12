
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/reachabilityanalysisintents` Documentation

The `reachabilityanalysisintents` SDK allows for interaction with Azure Resource Manager `network` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/reachabilityanalysisintents"
```


### Client Initialization

```go
client := reachabilityanalysisintents.NewReachabilityAnalysisIntentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReachabilityAnalysisIntentsClient.Create`

```go
ctx := context.TODO()
id := reachabilityanalysisintents.NewReachabilityAnalysisIntentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "verifierWorkspaceName", "reachabilityAnalysisIntentName")

payload := reachabilityanalysisintents.ReachabilityAnalysisIntent{
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


### Example Usage: `ReachabilityAnalysisIntentsClient.Get`

```go
ctx := context.TODO()
id := reachabilityanalysisintents.NewReachabilityAnalysisIntentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "verifierWorkspaceName", "reachabilityAnalysisIntentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReachabilityAnalysisIntentsClient.List`

```go
ctx := context.TODO()
id := reachabilityanalysisintents.NewVerifierWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "verifierWorkspaceName")

// alternatively `client.List(ctx, id, reachabilityanalysisintents.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, reachabilityanalysisintents.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
