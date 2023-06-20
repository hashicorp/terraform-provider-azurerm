
## `github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettings` Documentation

The `diagnosticsettings` SDK allows for interaction with the Azure Resource Manager Service `insights` (API Version `2021-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettings"
```


### Client Initialization

```go
client := diagnosticsettings.NewDiagnosticSettingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DiagnosticSettingsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := diagnosticsettings.NewScopedDiagnosticSettingID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "diagnosticSettingValue")

payload := diagnosticsettings.DiagnosticSettingsResource{
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


### Example Usage: `DiagnosticSettingsClient.Delete`

```go
ctx := context.TODO()
id := diagnosticsettings.NewScopedDiagnosticSettingID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "diagnosticSettingValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiagnosticSettingsClient.Get`

```go
ctx := context.TODO()
id := diagnosticsettings.NewScopedDiagnosticSettingID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "diagnosticSettingValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiagnosticSettingsClient.List`

```go
ctx := context.TODO()
id := diagnosticsettings.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
