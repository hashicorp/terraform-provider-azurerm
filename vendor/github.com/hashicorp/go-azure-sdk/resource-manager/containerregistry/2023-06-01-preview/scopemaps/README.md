
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/scopemaps` Documentation

The `scopemaps` SDK allows for interaction with Azure Resource Manager `containerregistry` (API Version `2023-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/scopemaps"
```


### Client Initialization

```go
client := scopemaps.NewScopeMapsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ScopeMapsClient.Create`

```go
ctx := context.TODO()
id := scopemaps.NewScopeMapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "scopeMapName")

payload := scopemaps.ScopeMap{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ScopeMapsClient.Delete`

```go
ctx := context.TODO()
id := scopemaps.NewScopeMapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "scopeMapName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ScopeMapsClient.Get`

```go
ctx := context.TODO()
id := scopemaps.NewScopeMapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "scopeMapName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScopeMapsClient.List`

```go
ctx := context.TODO()
id := scopemaps.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ScopeMapsClient.Update`

```go
ctx := context.TODO()
id := scopemaps.NewScopeMapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "scopeMapName")

payload := scopemaps.ScopeMapUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
