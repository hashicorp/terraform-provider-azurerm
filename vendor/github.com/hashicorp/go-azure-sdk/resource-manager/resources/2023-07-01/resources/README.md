
## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resources` Documentation

The `resources` SDK allows for interaction with Azure Resource Manager `resources` (API Version `2023-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resources"
```


### Client Initialization

```go
client := resources.NewResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ResourcesClient.CheckExistence`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.CheckExistence(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourcesClient.CheckExistenceById`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.CheckExistenceById(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourcesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := resources.GenericResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ResourcesClient.CreateOrUpdateById`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := resources.GenericResource{
	// ...
}


if err := client.CreateOrUpdateByIdThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ResourcesClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ResourcesClient.DeleteById`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

if err := client.DeleteByIdThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ResourcesClient.Get`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourcesClient.GetById`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.GetById(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourcesClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id, resources.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, resources.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourcesClient.MoveResources`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := resources.ResourcesMoveInfo{
	// ...
}


if err := client.MoveResourcesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ResourcesClient.Update`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := resources.GenericResource{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ResourcesClient.UpdateById`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := resources.GenericResource{
	// ...
}


if err := client.UpdateByIdThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ResourcesClient.ValidateMoveResources`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := resources.ResourcesMoveInfo{
	// ...
}


if err := client.ValidateMoveResourcesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
