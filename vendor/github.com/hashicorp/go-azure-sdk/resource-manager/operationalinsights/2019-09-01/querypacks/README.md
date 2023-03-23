
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


### Example Usage: `QueryPacksClient.QueryPacksCreateOrUpdate`

```go
ctx := context.TODO()
id := querypacks.NewQueryPackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue")

payload := querypacks.LogAnalyticsQueryPack{
	// ...
}


read, err := client.QueryPacksCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueryPacksClient.QueryPacksCreateOrUpdateWithoutName`

```go
ctx := context.TODO()
id := querypacks.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := querypacks.LogAnalyticsQueryPack{
	// ...
}


read, err := client.QueryPacksCreateOrUpdateWithoutName(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueryPacksClient.QueryPacksDelete`

```go
ctx := context.TODO()
id := querypacks.NewQueryPackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue")

read, err := client.QueryPacksDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueryPacksClient.QueryPacksGet`

```go
ctx := context.TODO()
id := querypacks.NewQueryPackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue")

read, err := client.QueryPacksGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueryPacksClient.QueryPacksList`

```go
ctx := context.TODO()
id := querypacks.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.QueryPacksList(ctx, id)` can be used to do batched pagination
items, err := client.QueryPacksListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `QueryPacksClient.QueryPacksListByResourceGroup`

```go
ctx := context.TODO()
id := querypacks.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.QueryPacksListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.QueryPacksListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `QueryPacksClient.QueryPacksUpdateTags`

```go
ctx := context.TODO()
id := querypacks.NewQueryPackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue")

payload := querypacks.TagsResource{
	// ...
}


read, err := client.QueryPacksUpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
