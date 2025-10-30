
## `github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/manageddashboards` Documentation

The `manageddashboards` SDK allows for interaction with Azure Resource Manager `dashboard` (API Version `2025-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/manageddashboards"
```


### Client Initialization

```go
client := manageddashboards.NewManagedDashboardsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedDashboardsClient.Create`

```go
ctx := context.TODO()
id := manageddashboards.NewDashboardID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dashboardName")

payload := manageddashboards.ManagedDashboard{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedDashboardsClient.DashboardsGet`

```go
ctx := context.TODO()
id := manageddashboards.NewDashboardID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dashboardName")

read, err := client.DashboardsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedDashboardsClient.DashboardsList`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.DashboardsList(ctx, id)` can be used to do batched pagination
items, err := client.DashboardsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedDashboardsClient.DashboardsListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.DashboardsListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.DashboardsListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedDashboardsClient.Delete`

```go
ctx := context.TODO()
id := manageddashboards.NewDashboardID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dashboardName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedDashboardsClient.Update`

```go
ctx := context.TODO()
id := manageddashboards.NewDashboardID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dashboardName")

payload := manageddashboards.ManagedDashboardUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
