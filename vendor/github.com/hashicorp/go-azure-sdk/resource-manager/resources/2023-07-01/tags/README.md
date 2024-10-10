
## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/tags` Documentation

The `tags` SDK allows for interaction with Azure Resource Manager `resources` (API Version `2023-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/tags"
```


### Client Initialization

```go
client := tags.NewTagsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TagsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := tags.NewTagNameID("12345678-1234-9876-4563-123456789012", "tagName")

read, err := client.CreateOrUpdate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagsClient.CreateOrUpdateAtScope`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := tags.TagsResource{
	// ...
}


if err := client.CreateOrUpdateAtScopeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TagsClient.CreateOrUpdateValue`

```go
ctx := context.TODO()
id := tags.NewTagValueID("12345678-1234-9876-4563-123456789012", "tagName", "tagValueName")

read, err := client.CreateOrUpdateValue(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagsClient.Delete`

```go
ctx := context.TODO()
id := tags.NewTagNameID("12345678-1234-9876-4563-123456789012", "tagName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagsClient.DeleteAtScope`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

if err := client.DeleteAtScopeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `TagsClient.DeleteValue`

```go
ctx := context.TODO()
id := tags.NewTagValueID("12345678-1234-9876-4563-123456789012", "tagName", "tagValueName")

read, err := client.DeleteValue(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagsClient.GetAtScope`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.GetAtScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TagsClient.UpdateAtScope`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := tags.TagsPatchResource{
	// ...
}


if err := client.UpdateAtScopeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
