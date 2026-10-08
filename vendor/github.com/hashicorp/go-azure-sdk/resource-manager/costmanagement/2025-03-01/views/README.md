
## `github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/views` Documentation

The `views` SDK allows for interaction with Azure Resource Manager `costmanagement` (API Version `2025-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/views"
```


### Client Initialization

```go
client := views.NewViewsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ViewsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := views.NewViewID("viewName")

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


### Example Usage: `ViewsClient.Delete`

```go
ctx := context.TODO()
id := views.NewViewID("viewName")

read, err := client.Delete(ctx, id)
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
id := views.NewViewID("viewName")

read, err := client.Get(ctx, id)
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
