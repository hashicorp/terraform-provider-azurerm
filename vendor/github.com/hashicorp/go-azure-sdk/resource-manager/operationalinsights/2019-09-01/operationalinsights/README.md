
## `github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/operationalinsights` Documentation

The `operationalinsights` SDK allows for interaction with the Azure Resource Manager Service `operationalinsights` (API Version `2019-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/operationalinsights"
```


### Client Initialization

```go
client := operationalinsights.NewOperationalInsightsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OperationalInsightsClient.QueriesDelete`

```go
ctx := context.TODO()
id := operationalinsights.NewQueriesID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue", "idValue")

read, err := client.QueriesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OperationalInsightsClient.QueriesGet`

```go
ctx := context.TODO()
id := operationalinsights.NewQueriesID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue", "idValue")

read, err := client.QueriesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OperationalInsightsClient.QueriesList`

```go
ctx := context.TODO()
id := operationalinsights.NewQueryPackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue")

// alternatively `client.QueriesList(ctx, id, operationalinsights.DefaultQueriesListOperationOptions())` can be used to do batched pagination
items, err := client.QueriesListComplete(ctx, id, operationalinsights.DefaultQueriesListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OperationalInsightsClient.QueriesPut`

```go
ctx := context.TODO()
id := operationalinsights.NewQueriesID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue", "idValue")

payload := operationalinsights.LogAnalyticsQueryPackQuery{
	// ...
}


read, err := client.QueriesPut(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OperationalInsightsClient.QueriesSearch`

```go
ctx := context.TODO()
id := operationalinsights.NewQueryPackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue")

payload := operationalinsights.LogAnalyticsQueryPackQuerySearchProperties{
	// ...
}


// alternatively `client.QueriesSearch(ctx, id, payload, operationalinsights.DefaultQueriesSearchOperationOptions())` can be used to do batched pagination
items, err := client.QueriesSearchComplete(ctx, id, payload, operationalinsights.DefaultQueriesSearchOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OperationalInsightsClient.QueriesUpdate`

```go
ctx := context.TODO()
id := operationalinsights.NewQueriesID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue", "idValue")

payload := operationalinsights.LogAnalyticsQueryPackQuery{
	// ...
}


read, err := client.QueriesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OperationalInsightsClient.QueryPacksCreateOrUpdate`

```go
ctx := context.TODO()
id := operationalinsights.NewQueryPackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue")

payload := operationalinsights.LogAnalyticsQueryPack{
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


### Example Usage: `OperationalInsightsClient.QueryPacksDelete`

```go
ctx := context.TODO()
id := operationalinsights.NewQueryPackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue")

read, err := client.QueryPacksDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OperationalInsightsClient.QueryPacksGet`

```go
ctx := context.TODO()
id := operationalinsights.NewQueryPackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue")

read, err := client.QueryPacksGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OperationalInsightsClient.QueryPacksList`

```go
ctx := context.TODO()
id := operationalinsights.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.QueryPacksList(ctx, id)` can be used to do batched pagination
items, err := client.QueryPacksListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OperationalInsightsClient.QueryPacksListByResourceGroup`

```go
ctx := context.TODO()
id := operationalinsights.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.QueryPacksListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.QueryPacksListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OperationalInsightsClient.QueryPacksUpdateTags`

```go
ctx := context.TODO()
id := operationalinsights.NewQueryPackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryPackValue")

payload := operationalinsights.TagsResource{
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
