
## `github.com/hashicorp/go-azure-sdk/resource-manager/security/2022-05-01/settings` Documentation

The `settings` SDK allows for interaction with Azure Resource Manager `security` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/security/2022-05-01/settings"
```


### Client Initialization

```go
client := settings.NewSettingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SettingsClient.Get`

```go
ctx := context.TODO()
id := settings.NewSettingID("12345678-1234-9876-4563-123456789012", "MCAS")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SettingsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SettingsClient.Update`

```go
ctx := context.TODO()
id := settings.NewSettingID("12345678-1234-9876-4563-123456789012", "MCAS")

payload := settings.Setting{
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
