
## `github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/dashboard` Documentation

The `dashboard` SDK allows for interaction with the Azure Resource Manager Service `portal` (API Version `2019-01-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/dashboard"
```


### Client Initialization

```go
client := dashboard.NewDashboardClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DashboardClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dashboard.NewDashboardID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dashboardValue")

payload := dashboard.Dashboard{
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


### Example Usage: `DashboardClient.Delete`

```go
ctx := context.TODO()
id := dashboard.NewDashboardID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dashboardValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DashboardClient.Get`

```go
ctx := context.TODO()
id := dashboard.NewDashboardID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dashboardValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DashboardClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := dashboard.NewResourceGroupID()

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DashboardClient.ListBySubscription`

```go
ctx := context.TODO()
id := dashboard.NewSubscriptionID()

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DashboardClient.Update`

```go
ctx := context.TODO()
id := dashboard.NewDashboardID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dashboardValue")

payload := dashboard.PatchableDashboard{
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
