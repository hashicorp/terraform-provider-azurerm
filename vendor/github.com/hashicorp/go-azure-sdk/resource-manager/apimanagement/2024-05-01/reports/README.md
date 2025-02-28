
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/reports` Documentation

The `reports` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/reports"
```


### Client Initialization

```go
client := reports.NewReportsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReportsClient.ListByApi`

```go
ctx := context.TODO()
id := reports.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByApi(ctx, id, reports.DefaultListByApiOperationOptions())` can be used to do batched pagination
items, err := client.ListByApiComplete(ctx, id, reports.DefaultListByApiOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReportsClient.ListByGeo`

```go
ctx := context.TODO()
id := reports.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByGeo(ctx, id, reports.DefaultListByGeoOperationOptions())` can be used to do batched pagination
items, err := client.ListByGeoComplete(ctx, id, reports.DefaultListByGeoOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReportsClient.ListByOperation`

```go
ctx := context.TODO()
id := reports.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByOperation(ctx, id, reports.DefaultListByOperationOperationOptions())` can be used to do batched pagination
items, err := client.ListByOperationComplete(ctx, id, reports.DefaultListByOperationOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReportsClient.ListByProduct`

```go
ctx := context.TODO()
id := reports.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByProduct(ctx, id, reports.DefaultListByProductOperationOptions())` can be used to do batched pagination
items, err := client.ListByProductComplete(ctx, id, reports.DefaultListByProductOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReportsClient.ListByRequest`

```go
ctx := context.TODO()
id := reports.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

read, err := client.ListByRequest(ctx, id, reports.DefaultListByRequestOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReportsClient.ListBySubscription`

```go
ctx := context.TODO()
id := reports.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListBySubscription(ctx, id, reports.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, reports.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReportsClient.ListByTime`

```go
ctx := context.TODO()
id := reports.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByTime(ctx, id, reports.DefaultListByTimeOperationOptions())` can be used to do batched pagination
items, err := client.ListByTimeComplete(ctx, id, reports.DefaultListByTimeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReportsClient.ListByUser`

```go
ctx := context.TODO()
id := reports.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByUser(ctx, id, reports.DefaultListByUserOperationOptions())` can be used to do batched pagination
items, err := client.ListByUserComplete(ctx, id, reports.DefaultListByUserOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
