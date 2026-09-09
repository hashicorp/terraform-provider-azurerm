
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/routetables` Documentation

The `routetables` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/routetables"
```


### Client Initialization

```go
client := routetables.NewRouteTablesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RouteTablesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := routetables.NewRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeTableName")

payload := routetables.RouteTable{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RouteTablesClient.Delete`

```go
ctx := context.TODO()
id := routetables.NewRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeTableName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RouteTablesClient.Get`

```go
ctx := context.TODO()
id := routetables.NewRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeTableName")

read, err := client.Get(ctx, id, routetables.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RouteTablesClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RouteTablesClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RouteTablesClient.UpdateTags`

```go
ctx := context.TODO()
id := routetables.NewRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeTableName")

payload := routetables.TagsObject{
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
