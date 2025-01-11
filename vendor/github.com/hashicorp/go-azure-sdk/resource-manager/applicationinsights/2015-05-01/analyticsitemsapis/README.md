
## `github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/analyticsitemsapis` Documentation

The `analyticsitemsapis` SDK allows for interaction with Azure Resource Manager `applicationinsights` (API Version `2015-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/analyticsitemsapis"
```


### Client Initialization

```go
client := analyticsitemsapis.NewAnalyticsItemsAPIsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AnalyticsItemsAPIsClient.AnalyticsItemsDelete`

```go
ctx := context.TODO()
id := analyticsitemsapis.NewProviderComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.AnalyticsItemsDelete(ctx, id, analyticsitemsapis.DefaultAnalyticsItemsDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AnalyticsItemsAPIsClient.AnalyticsItemsGet`

```go
ctx := context.TODO()
id := analyticsitemsapis.NewProviderComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.AnalyticsItemsGet(ctx, id, analyticsitemsapis.DefaultAnalyticsItemsGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AnalyticsItemsAPIsClient.AnalyticsItemsList`

```go
ctx := context.TODO()
id := analyticsitemsapis.NewProviderComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.AnalyticsItemsList(ctx, id, analyticsitemsapis.DefaultAnalyticsItemsListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AnalyticsItemsAPIsClient.AnalyticsItemsPut`

```go
ctx := context.TODO()
id := analyticsitemsapis.NewProviderComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := analyticsitemsapis.ApplicationInsightsComponentAnalyticsItem{
	// ...
}


read, err := client.AnalyticsItemsPut(ctx, id, payload, analyticsitemsapis.DefaultAnalyticsItemsPutOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
