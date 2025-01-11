
## `github.com/hashicorp/go-azure-sdk/resource-manager/insights/2019-10-17-preview/privatelinkscopedresources` Documentation

The `privatelinkscopedresources` SDK allows for interaction with Azure Resource Manager `insights` (API Version `2019-10-17-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2019-10-17-preview/privatelinkscopedresources"
```


### Client Initialization

```go
client := privatelinkscopedresources.NewPrivateLinkScopedResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateLinkScopedResourcesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := privatelinkscopedresources.NewScopedResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkScopeName", "scopedResourceName")

payload := privatelinkscopedresources.ScopedResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateLinkScopedResourcesClient.Delete`

```go
ctx := context.TODO()
id := privatelinkscopedresources.NewScopedResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkScopeName", "scopedResourceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateLinkScopedResourcesClient.Get`

```go
ctx := context.TODO()
id := privatelinkscopedresources.NewScopedResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkScopeName", "scopedResourceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkScopedResourcesClient.ListByPrivateLinkScope`

```go
ctx := context.TODO()
id := privatelinkscopedresources.NewPrivateLinkScopeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkScopeName")

// alternatively `client.ListByPrivateLinkScope(ctx, id)` can be used to do batched pagination
items, err := client.ListByPrivateLinkScopeComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
