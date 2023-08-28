
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/routefilters` Documentation

The `routefilters` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/routefilters"
```


### Client Initialization

```go
client := routefilters.NewRouteFiltersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RouteFiltersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := routefilters.NewRouteFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeFilterValue")

payload := routefilters.RouteFilter{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RouteFiltersClient.Delete`

```go
ctx := context.TODO()
id := routefilters.NewRouteFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeFilterValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RouteFiltersClient.Get`

```go
ctx := context.TODO()
id := routefilters.NewRouteFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeFilterValue")

read, err := client.Get(ctx, id, routefilters.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RouteFiltersClient.List`

```go
ctx := context.TODO()
id := routefilters.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RouteFiltersClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := routefilters.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RouteFiltersClient.UpdateTags`

```go
ctx := context.TODO()
id := routefilters.NewRouteFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeFilterValue")

payload := routefilters.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
