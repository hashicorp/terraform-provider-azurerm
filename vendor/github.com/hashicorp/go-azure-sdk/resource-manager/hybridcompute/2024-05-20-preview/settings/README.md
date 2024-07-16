
## `github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-05-20-preview/settings` Documentation

The `settings` SDK allows for interaction with the Azure Resource Manager Service `hybridcompute` (API Version `2024-05-20-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-05-20-preview/settings"
```


### Client Initialization

```go
client := settings.NewSettingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SettingsClient.Get`

```go
ctx := context.TODO()
id := settings.NewScopedSettingID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "settingValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SettingsClient.Patch`

```go
ctx := context.TODO()
id := settings.NewScopedSettingID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "settingValue")

payload := settings.Settings{
	// ...
}


read, err := client.Patch(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SettingsClient.Update`

```go
ctx := context.TODO()
id := settings.NewScopedSettingID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "settingValue")

payload := settings.Settings{
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
