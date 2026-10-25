
## `github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/viewoperationgroup` Documentation

The `viewoperationgroup` SDK allows for interaction with Azure Resource Manager `costmanagement` (API Version `2025-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/viewoperationgroup"
```


### Client Initialization

```go
client := viewoperationgroup.NewViewOperationGroupClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ViewOperationGroupClient.ViewsCreateOrUpdateByScope`

```go
ctx := context.TODO()
id := viewoperationgroup.NewScopedViewID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "viewName")

payload := viewoperationgroup.View{
	// ...
}


read, err := client.ViewsCreateOrUpdateByScope(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ViewOperationGroupClient.ViewsDeleteByScope`

```go
ctx := context.TODO()
id := viewoperationgroup.NewScopedViewID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "viewName")

read, err := client.ViewsDeleteByScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ViewOperationGroupClient.ViewsGetByScope`

```go
ctx := context.TODO()
id := viewoperationgroup.NewScopedViewID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "viewName")

read, err := client.ViewsGetByScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ViewOperationGroupClient.ViewsListByScope`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ViewsListByScope(ctx, id)` can be used to do batched pagination
items, err := client.ViewsListByScopeComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
