
## `github.com/hashicorp/go-azure-sdk/resource-manager/resourcegraph/2024-04-01/graphquery` Documentation

The `graphquery` SDK allows for interaction with Azure Resource Manager `resourcegraph` (API Version `2024-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resourcegraph/2024-04-01/graphquery"
```


### Client Initialization

```go
client := graphquery.NewGraphQueryClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GraphQueryClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := graphquery.NewQueryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryName")

payload := graphquery.GraphQueryResource{
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


### Example Usage: `GraphQueryClient.Delete`

```go
ctx := context.TODO()
id := graphquery.NewQueryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GraphQueryClient.Get`

```go
ctx := context.TODO()
id := graphquery.NewQueryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GraphQueryClient.List`

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


### Example Usage: `GraphQueryClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
