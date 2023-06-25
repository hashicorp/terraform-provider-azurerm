
## `github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/synchronizationsetting` Documentation

The `synchronizationsetting` SDK allows for interaction with the Azure Resource Manager Service `datashare` (API Version `2019-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/synchronizationsetting"
```


### Client Initialization

```go
client := synchronizationsetting.NewSynchronizationSettingClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SynchronizationSettingClient.Create`

```go
ctx := context.TODO()
id := synchronizationsetting.NewSynchronizationSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "shareValue", "synchronizationSettingValue")

payload := synchronizationsetting.SynchronizationSetting{
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


### Example Usage: `SynchronizationSettingClient.Delete`

```go
ctx := context.TODO()
id := synchronizationsetting.NewSynchronizationSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "shareValue", "synchronizationSettingValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SynchronizationSettingClient.Get`

```go
ctx := context.TODO()
id := synchronizationsetting.NewSynchronizationSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "shareValue", "synchronizationSettingValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SynchronizationSettingClient.ListByShare`

```go
ctx := context.TODO()
id := synchronizationsetting.NewShareID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "shareValue")

// alternatively `client.ListByShare(ctx, id)` can be used to do batched pagination
items, err := client.ListByShareComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
