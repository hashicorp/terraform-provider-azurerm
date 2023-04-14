
## `github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2022-10-01/views` Documentation

The `views` SDK allows for interaction with the Azure Resource Manager Service `costmanagement` (API Version `2022-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2022-10-01/views"
```


### Client Initialization

```go
client := views.NewViewsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ViewsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := views.NewViewID("viewValue")

payload := views.View{
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


### Example Usage: `ViewsClient.CreateOrUpdateByScope`

```go
ctx := context.TODO()
id := views.NewScopedViewID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "viewValue")

payload := views.View{
	// ...
}


read, err := client.CreateOrUpdateByScope(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ViewsClient.Delete`

```go
ctx := context.TODO()
id := views.NewViewID("viewValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ViewsClient.DeleteByScope`

```go
ctx := context.TODO()
id := views.NewScopedViewID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "viewValue")

read, err := client.DeleteByScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ViewsClient.Get`

```go
ctx := context.TODO()
id := views.NewViewID("viewValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ViewsClient.GetByScope`

```go
ctx := context.TODO()
id := views.NewScopedViewID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "viewValue")

read, err := client.GetByScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ViewsClient.List`

```go
ctx := context.TODO()


// alternatively `client.List(ctx)` can be used to do batched pagination
items, err := client.ListComplete(ctx)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ViewsClient.ListByScope`

```go
ctx := context.TODO()
id := views.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ListByScope(ctx, id)` can be used to do batched pagination
items, err := client.ListByScopeComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
