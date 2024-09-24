
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/scopeconnections` Documentation

The `scopeconnections` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/scopeconnections"
```


### Client Initialization

```go
client := scopeconnections.NewScopeConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ScopeConnectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := scopeconnections.NewScopeConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "scopeConnectionName")

payload := scopeconnections.ScopeConnection{
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


### Example Usage: `ScopeConnectionsClient.Delete`

```go
ctx := context.TODO()
id := scopeconnections.NewScopeConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "scopeConnectionName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScopeConnectionsClient.Get`

```go
ctx := context.TODO()
id := scopeconnections.NewScopeConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "scopeConnectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScopeConnectionsClient.List`

```go
ctx := context.TODO()
id := scopeconnections.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName")

// alternatively `client.List(ctx, id, scopeconnections.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, scopeconnections.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
