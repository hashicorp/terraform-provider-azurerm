
## `github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitors` Documentation

The `monitors` SDK allows for interaction with Azure Resource Manager `newrelic` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitors"
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

payload := monitors.NewRelicMonitorResource{
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

if err := client.DeleteThenPoll(ctx, id, monitors.DefaultDeleteOperationOptions()); err != nil {
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


### Example Usage: `MonitorsClient.GetMetricRules`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := monitors.MetricsRequest{
	// ...
}


read, err := client.GetMetricRules(ctx, id, payload)
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

payload := monitors.MetricsStatusRequest{
	// ...
}


read, err := client.GetMetricStatus(ctx, id, payload)
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

payload := monitors.AppServicesGetRequest{
	// ...
}


// alternatively `client.ListAppServices(ctx, id, payload)` can be used to do batched pagination
items, err := client.ListAppServicesComplete(ctx, id, payload)
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


### Example Usage: `MonitorsClient.ListBySubscription`

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


### Example Usage: `MonitorsClient.ListHosts`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := monitors.HostsGetRequest{
	// ...
}


// alternatively `client.ListHosts(ctx, id, payload)` can be used to do batched pagination
items, err := client.ListHostsComplete(ctx, id, payload)
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


### Example Usage: `MonitorsClient.SwitchBilling`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := monitors.SwitchBillingRequest{
	// ...
}


read, err := client.SwitchBilling(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MonitorsClient.Update`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := monitors.NewRelicMonitorResourceUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `MonitorsClient.VMHostPayload`

```go
ctx := context.TODO()
id := monitors.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.VMHostPayload(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
