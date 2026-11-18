
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/connectionanalyzers` Documentation

The `connectionanalyzers` SDK allows for interaction with Azure Resource Manager `network` (API Version `2026-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/connectionanalyzers"
```


### Client Initialization

```go
client := connectionanalyzers.NewConnectionAnalyzersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConnectionAnalyzersClient.NetworkWatchersConnectionAnalyzersCreate`

```go
ctx := context.TODO()
id := connectionanalyzers.NewConnectionAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "connectionAnalyzerName")

payload := connectionanalyzers.ConnectionAnalyzer{
	// ...
}


if err := client.NetworkWatchersConnectionAnalyzersCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ConnectionAnalyzersClient.NetworkWatchersConnectionAnalyzersDelete`

```go
ctx := context.TODO()
id := connectionanalyzers.NewConnectionAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "connectionAnalyzerName")

if err := client.NetworkWatchersConnectionAnalyzersDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ConnectionAnalyzersClient.NetworkWatchersConnectionAnalyzersGet`

```go
ctx := context.TODO()
id := connectionanalyzers.NewConnectionAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "connectionAnalyzerName")

read, err := client.NetworkWatchersConnectionAnalyzersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectionAnalyzersClient.NetworkWatchersConnectionAnalyzersList`

```go
ctx := context.TODO()
id := connectionanalyzers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

// alternatively `client.NetworkWatchersConnectionAnalyzersList(ctx, id)` can be used to do batched pagination
items, err := client.NetworkWatchersConnectionAnalyzersListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ConnectionAnalyzersClient.NetworkWatchersConnectionAnalyzersQuery`

```go
ctx := context.TODO()
id := connectionanalyzers.NewConnectionAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "connectionAnalyzerName")

if err := client.NetworkWatchersConnectionAnalyzersQueryThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ConnectionAnalyzersClient.NetworkWatchersConnectionAnalyzersUpdateTags`

```go
ctx := context.TODO()
id := connectionanalyzers.NewConnectionAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "connectionAnalyzerName")

payload := connectionanalyzers.TagsObject{
	// ...
}


read, err := client.NetworkWatchersConnectionAnalyzersUpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
