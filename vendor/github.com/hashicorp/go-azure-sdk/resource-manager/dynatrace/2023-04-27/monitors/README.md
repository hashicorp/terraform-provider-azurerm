
## `github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors` Documentation

The `monitors` SDK allows for interaction with Azure Resource Manager `dynatrace` (API Version `2023-04-27`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
```


### Client Initialization

```go
client := monitors.NewMonitorsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MonitorsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := monitors.MonitorResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `MonitorsClient.Delete`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `MonitorsClient.Get`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MonitorsClient.GetMarketplaceSaaSResourceDetails`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := monitors.MarketplaceSaaSResourceDetailsRequest{
	// ...
}


read, err := client.GetMarketplaceSaaSResourceDetails(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MonitorsClient.GetMetricStatus`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.GetMetricStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MonitorsClient.GetSSODetails`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := monitors.SSODetailsRequest{
	// ...
}


read, err := client.GetSSODetails(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MonitorsClient.GetVMHostPayload`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.GetVMHostPayload(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MonitorsClient.ListAppServices`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.ListAppServices(ctx, id)` can be used to do batched pagination
items, err := client.ListAppServicesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MonitorsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MonitorsClient.ListBySubscriptionId`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscriptionId(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionIdComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MonitorsClient.ListHosts`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.ListHosts(ctx, id)` can be used to do batched pagination
items, err := client.ListHostsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MonitorsClient.ListLinkableEnvironments`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := monitors.LinkableEnvironmentRequest{
	// ...
}


// alternatively `client.ListLinkableEnvironments(ctx, id, payload)` can be used to do batched pagination
items, err := client.ListLinkableEnvironmentsComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MonitorsClient.ListMonitoredResources`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.ListMonitoredResources(ctx, id)` can be used to do batched pagination
items, err := client.ListMonitoredResourcesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MonitorsClient.Update`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := monitors.MonitorResourceUpdate{
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
