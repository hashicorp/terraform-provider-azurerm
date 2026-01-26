
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/settings` Documentation

The `settings` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/settings"
```


### Client Initialization

```go
client := settings.NewSettingsClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `SettingsClient.GetSetting`

```go
ctx := context.TODO()
id := settings.NewSettingID("settingName")

read, err := client.GetSetting(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SettingsClient.GetSettings`

```go
ctx := context.TODO()


read, err := client.GetSettings(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SettingsClient.UpdateSetting`

```go
ctx := context.TODO()
id := settings.NewSettingID("settingName")

payload := settings.UpdateSettingRequest{
	// ...
}


read, err := client.UpdateSetting(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
