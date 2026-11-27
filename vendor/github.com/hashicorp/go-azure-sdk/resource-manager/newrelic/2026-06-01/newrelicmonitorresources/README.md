
## `github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2026-06-01/newrelicmonitorresources` Documentation

The `newrelicmonitorresources` SDK allows for interaction with Azure Resource Manager `newrelic` (API Version `2026-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2026-06-01/newrelicmonitorresources"
```


### Client Initialization

```go
client := newrelicmonitorresources.NewNewRelicMonitorResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NewRelicMonitorResourcesClient.BillingInfoGet`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.BillingInfoGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NewRelicMonitorResourcesClient.ConnectedPartnerResourcesList`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")
var payload string

// alternatively `client.ConnectedPartnerResourcesList(ctx, id, payload)` can be used to do batched pagination
items, err := client.ConnectedPartnerResourcesListComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsCreateOrUpdate`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := newrelicmonitorresources.NewRelicMonitorResource{
	// ...
}


if err := client.MonitorsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsDelete`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

if err := client.MonitorsDeleteThenPoll(ctx, id, newrelicmonitorresources.DefaultMonitorsDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsGet`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.MonitorsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsGetMetricRules`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := newrelicmonitorresources.MetricsRequest{
	// ...
}


read, err := client.MonitorsGetMetricRules(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsGetMetricStatus`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := newrelicmonitorresources.MetricsStatusRequest{
	// ...
}


read, err := client.MonitorsGetMetricStatus(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsLatestLinkedSaaS`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.MonitorsLatestLinkedSaaS(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsLinkSaaS`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := newrelicmonitorresources.SaaSData{
	// ...
}


if err := client.MonitorsLinkSaaSThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsListAppServices`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := newrelicmonitorresources.AppServicesGetRequest{
	// ...
}


// alternatively `client.MonitorsListAppServices(ctx, id, payload)` can be used to do batched pagination
items, err := client.MonitorsListAppServicesComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.MonitorsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.MonitorsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.MonitorsListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.MonitorsListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsListHosts`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := newrelicmonitorresources.HostsGetRequest{
	// ...
}


// alternatively `client.MonitorsListHosts(ctx, id, payload)` can be used to do batched pagination
items, err := client.MonitorsListHostsComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsListLinkedResources`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.MonitorsListLinkedResources(ctx, id)` can be used to do batched pagination
items, err := client.MonitorsListLinkedResourcesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsListMonitoredResources`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.MonitorsListMonitoredResources(ctx, id)` can be used to do batched pagination
items, err := client.MonitorsListMonitoredResourcesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsRefreshIngestionKey`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.MonitorsRefreshIngestionKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsResubscribe`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := newrelicmonitorresources.ResubscribeProperties{
	// ...
}


if err := client.MonitorsResubscribeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsSwitchBilling`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := newrelicmonitorresources.SwitchBillingRequest{
	// ...
}


read, err := client.MonitorsSwitchBilling(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsUpdate`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := newrelicmonitorresources.NewRelicMonitorResourceUpdate{
	// ...
}


if err := client.MonitorsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NewRelicMonitorResourcesClient.MonitorsVMHostPayload`

```go
ctx := context.TODO()
id := newrelicmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.MonitorsVMHostPayload(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
