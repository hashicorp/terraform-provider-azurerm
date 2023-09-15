
## `github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/querypacks` Documentation

The `querypacks` SDK allows for interaction with the Azure Resource Manager Service `operationalinsights` (API Version `2019-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/querypacks"
```


### Client Initialization

```go
client := querypacks.NewQueryPacksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `QueryPacksClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := querypacks.NewQueryPackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue")

payload := querypacks.LogAnalyticsQueryPack{
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


### Example Usage: `QueryPacksClient.CreateOrUpdateWithoutName`

```go
ctx := context.TODO()
id := querypacks.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := querypacks.LogAnalyticsQueryPack{
	// ...
}


read, err := client.CreateOrUpdateWithoutName(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueryPacksClient.Delete`

```go
ctx := context.TODO()
id := querypacks.NewQueryPackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueryPacksClient.Get`

```go
ctx := context.TODO()
id := querypacks.NewQueryPackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueryPacksClient.List`

```go
ctx := context.TODO()
id := querypacks.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `QueryPacksClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := querypacks.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `QueryPacksClient.UpdateTags`

```go
ctx := context.TODO()
id := querypacks.NewQueryPackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue")

payload := querypacks.TagsResource{
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
