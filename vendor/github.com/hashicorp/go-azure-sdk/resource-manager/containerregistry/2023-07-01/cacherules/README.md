
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/cacherules` Documentation

The `cacherules` SDK allows for interaction with Azure Resource Manager `containerregistry` (API Version `2023-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/cacherules"
```


### Client Initialization

```go
client := cacherules.NewCacheRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CacheRulesClient.Create`

```go
ctx := context.TODO()
id := cacherules.NewCacheRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "cacheRuleName")

payload := cacherules.CacheRule{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CacheRulesClient.Delete`

```go
ctx := context.TODO()
id := cacherules.NewCacheRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "cacheRuleName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CacheRulesClient.Get`

```go
ctx := context.TODO()
id := cacherules.NewCacheRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "cacheRuleName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CacheRulesClient.List`

```go
ctx := context.TODO()
id := cacherules.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CacheRulesClient.Update`

```go
ctx := context.TODO()
id := cacherules.NewCacheRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "cacheRuleName")

payload := cacherules.CacheRuleUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
