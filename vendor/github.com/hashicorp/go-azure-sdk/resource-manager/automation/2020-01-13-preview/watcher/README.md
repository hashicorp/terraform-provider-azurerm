
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/watcher` Documentation

The `watcher` SDK allows for interaction with the Azure Resource Manager Service `automation` (API Version `2020-01-13-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/watcher"
```


### Client Initialization

```go
client := watcher.NewWatcherClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `WatcherClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := watcher.NewWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "watcherValue")

payload := watcher.Watcher{
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


### Example Usage: `WatcherClient.Delete`

```go
ctx := context.TODO()
id := watcher.NewWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "watcherValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WatcherClient.Get`

```go
ctx := context.TODO()
id := watcher.NewWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "watcherValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WatcherClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := watcher.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue")

// alternatively `client.ListByAutomationAccount(ctx, id, watcher.DefaultListByAutomationAccountOperationOptions())` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id, watcher.DefaultListByAutomationAccountOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WatcherClient.Start`

```go
ctx := context.TODO()
id := watcher.NewWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "watcherValue")

read, err := client.Start(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WatcherClient.Stop`

```go
ctx := context.TODO()
id := watcher.NewWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "watcherValue")

read, err := client.Stop(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WatcherClient.Update`

```go
ctx := context.TODO()
id := watcher.NewWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "watcherValue")

payload := watcher.WatcherUpdateParameters{
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
