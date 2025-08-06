
## `github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/watchlists` Documentation

The `watchlists` SDK allows for interaction with Azure Resource Manager `securityinsights` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/watchlists"
```


### Client Initialization

```go
client := watchlists.NewWatchlistsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `WatchlistsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := watchlists.NewWatchlistID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "watchlistAlias")

payload := watchlists.Watchlist{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WatchlistsClient.Delete`

```go
ctx := context.TODO()
id := watchlists.NewWatchlistID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "watchlistAlias")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WatchlistsClient.Get`

```go
ctx := context.TODO()
id := watchlists.NewWatchlistID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "watchlistAlias")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WatchlistsClient.List`

```go
ctx := context.TODO()
id := watchlists.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
