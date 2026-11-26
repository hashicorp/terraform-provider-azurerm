
## `github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/monitoredsubscriptions` Documentation

The `monitoredsubscriptions` SDK allows for interaction with Azure Resource Manager `datadog` (API Version `2025-06-11`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/monitoredsubscriptions"
```


### Client Initialization

```go
client := monitoredsubscriptions.NewMonitoredSubscriptionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MonitoredSubscriptionsClient.CreateorUpdate`

```go
ctx := context.TODO()
id := monitoredsubscriptions.NewMonitoredSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName", "monitoredSubscriptionName")

payload := monitoredsubscriptions.MonitoredSubscriptionProperties{
	// ...
}


if err := client.CreateorUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `MonitoredSubscriptionsClient.Delete`

```go
ctx := context.TODO()
id := monitoredsubscriptions.NewMonitoredSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName", "monitoredSubscriptionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `MonitoredSubscriptionsClient.Get`

```go
ctx := context.TODO()
id := monitoredsubscriptions.NewMonitoredSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName", "monitoredSubscriptionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MonitoredSubscriptionsClient.List`

```go
ctx := context.TODO()
id := monitoredsubscriptions.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MonitoredSubscriptionsClient.Update`

```go
ctx := context.TODO()
id := monitoredsubscriptions.NewMonitoredSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName", "monitoredSubscriptionName")

payload := monitoredsubscriptions.MonitoredSubscriptionProperties{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
