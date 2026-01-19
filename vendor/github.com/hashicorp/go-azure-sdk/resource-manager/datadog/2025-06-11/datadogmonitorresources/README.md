
## `github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogmonitorresources` Documentation

The `datadogmonitorresources` SDK allows for interaction with Azure Resource Manager `datadog` (API Version `2025-06-11`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogmonitorresources"
```


### Client Initialization

```go
client := datadogmonitorresources.NewDatadogMonitorResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DatadogMonitorResourcesClient.BillingInfoGet`

```go
ctx := context.TODO()
id := datadogmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.BillingInfoGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatadogMonitorResourcesClient.MonitorsCreate`

```go
ctx := context.TODO()
id := datadogmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := datadogmonitorresources.DatadogMonitorResource{
	// ...
}


if err := client.MonitorsCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DatadogMonitorResourcesClient.MonitorsDelete`

```go
ctx := context.TODO()
id := datadogmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

if err := client.MonitorsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DatadogMonitorResourcesClient.MonitorsGet`

```go
ctx := context.TODO()
id := datadogmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.MonitorsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatadogMonitorResourcesClient.MonitorsGetDefaultKey`

```go
ctx := context.TODO()
id := datadogmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.MonitorsGetDefaultKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatadogMonitorResourcesClient.MonitorsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.MonitorsList(ctx, id)` can be used to do batched pagination
items, err := client.MonitorsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DatadogMonitorResourcesClient.MonitorsListApiKeys`

```go
ctx := context.TODO()
id := datadogmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.MonitorsListApiKeys(ctx, id)` can be used to do batched pagination
items, err := client.MonitorsListApiKeysComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DatadogMonitorResourcesClient.MonitorsListByResourceGroup`

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


### Example Usage: `DatadogMonitorResourcesClient.MonitorsListHosts`

```go
ctx := context.TODO()
id := datadogmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.MonitorsListHosts(ctx, id)` can be used to do batched pagination
items, err := client.MonitorsListHostsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DatadogMonitorResourcesClient.MonitorsListLinkedResources`

```go
ctx := context.TODO()
id := datadogmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.MonitorsListLinkedResources(ctx, id)` can be used to do batched pagination
items, err := client.MonitorsListLinkedResourcesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DatadogMonitorResourcesClient.MonitorsListMonitoredResources`

```go
ctx := context.TODO()
id := datadogmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.MonitorsListMonitoredResources(ctx, id)` can be used to do batched pagination
items, err := client.MonitorsListMonitoredResourcesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DatadogMonitorResourcesClient.MonitorsRefreshSetPasswordLink`

```go
ctx := context.TODO()
id := datadogmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.MonitorsRefreshSetPasswordLink(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatadogMonitorResourcesClient.MonitorsSetDefaultKey`

```go
ctx := context.TODO()
id := datadogmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := datadogmonitorresources.DatadogApiKey{
	// ...
}


read, err := client.MonitorsSetDefaultKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatadogMonitorResourcesClient.MonitorsUpdate`

```go
ctx := context.TODO()
id := datadogmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := datadogmonitorresources.DatadogMonitorResourceUpdateParameters{
	// ...
}


if err := client.MonitorsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DatadogMonitorResourcesClient.OrganizationsResubscribe`

```go
ctx := context.TODO()
id := datadogmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := datadogmonitorresources.ResubscribeProperties{
	// ...
}


if err := client.OrganizationsResubscribeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
