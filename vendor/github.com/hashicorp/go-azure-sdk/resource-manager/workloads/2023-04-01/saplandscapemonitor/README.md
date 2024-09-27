
## `github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/saplandscapemonitor` Documentation

The `saplandscapemonitor` SDK allows for interaction with Azure Resource Manager `workloads` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/saplandscapemonitor"
```


### Client Initialization

```go
client := saplandscapemonitor.NewSapLandscapeMonitorClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SapLandscapeMonitorClient.Create`

```go
ctx := context.TODO()
id := saplandscapemonitor.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := saplandscapemonitor.SapLandscapeMonitor{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SapLandscapeMonitorClient.Delete`

```go
ctx := context.TODO()
id := saplandscapemonitor.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SapLandscapeMonitorClient.Get`

```go
ctx := context.TODO()
id := saplandscapemonitor.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SapLandscapeMonitorClient.List`

```go
ctx := context.TODO()
id := saplandscapemonitor.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SapLandscapeMonitorClient.Update`

```go
ctx := context.TODO()
id := saplandscapemonitor.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := saplandscapemonitor.SapLandscapeMonitor{
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
