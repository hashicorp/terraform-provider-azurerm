
## `github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-03-01/metricalerts` Documentation

The `metricalerts` SDK allows for interaction with Azure Resource Manager `insights` (API Version `2018-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-03-01/metricalerts"
```


### Client Initialization

```go
client := metricalerts.NewMetricAlertsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MetricAlertsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := metricalerts.NewMetricAlertID("12345678-1234-9876-4563-123456789012", "example-resource-group", "metricAlertName")

payload := metricalerts.MetricAlertResource{
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


### Example Usage: `MetricAlertsClient.Delete`

```go
ctx := context.TODO()
id := metricalerts.NewMetricAlertID("12345678-1234-9876-4563-123456789012", "example-resource-group", "metricAlertName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MetricAlertsClient.Get`

```go
ctx := context.TODO()
id := metricalerts.NewMetricAlertID("12345678-1234-9876-4563-123456789012", "example-resource-group", "metricAlertName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MetricAlertsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.ListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MetricAlertsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListBySubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MetricAlertsClient.Update`

```go
ctx := context.TODO()
id := metricalerts.NewMetricAlertID("12345678-1234-9876-4563-123456789012", "example-resource-group", "metricAlertName")

payload := metricalerts.MetricAlertResourcePatch{
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
