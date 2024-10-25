
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/connectionmonitors` Documentation

The `connectionmonitors` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/connectionmonitors"
```


### Client Initialization

```go
client := connectionmonitors.NewConnectionMonitorsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConnectionMonitorsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := connectionmonitors.NewConnectionMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "connectionMonitorName")

payload := connectionmonitors.ConnectionMonitor{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, connectionmonitors.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ConnectionMonitorsClient.Delete`

```go
ctx := context.TODO()
id := connectionmonitors.NewConnectionMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "connectionMonitorName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ConnectionMonitorsClient.Get`

```go
ctx := context.TODO()
id := connectionmonitors.NewConnectionMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "connectionMonitorName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectionMonitorsClient.List`

```go
ctx := context.TODO()
id := connectionmonitors.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectionMonitorsClient.Query`

```go
ctx := context.TODO()
id := connectionmonitors.NewConnectionMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "connectionMonitorName")

if err := client.QueryThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ConnectionMonitorsClient.Start`

```go
ctx := context.TODO()
id := connectionmonitors.NewConnectionMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "connectionMonitorName")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ConnectionMonitorsClient.Stop`

```go
ctx := context.TODO()
id := connectionmonitors.NewConnectionMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "connectionMonitorName")

if err := client.StopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ConnectionMonitorsClient.UpdateTags`

```go
ctx := context.TODO()
id := connectionmonitors.NewConnectionMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "connectionMonitorName")

payload := connectionmonitors.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
