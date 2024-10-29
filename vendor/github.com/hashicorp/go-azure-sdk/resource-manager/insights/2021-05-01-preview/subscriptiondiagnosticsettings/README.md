
## `github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/subscriptiondiagnosticsettings` Documentation

The `subscriptiondiagnosticsettings` SDK allows for interaction with Azure Resource Manager `insights` (API Version `2021-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/subscriptiondiagnosticsettings"
```


### Client Initialization

```go
client := subscriptiondiagnosticsettings.NewSubscriptionDiagnosticSettingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SubscriptionDiagnosticSettingsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := subscriptiondiagnosticsettings.NewDiagnosticSettingID("12345678-1234-9876-4563-123456789012", "diagnosticSettingName")

payload := subscriptiondiagnosticsettings.SubscriptionDiagnosticSettingsResource{
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


### Example Usage: `SubscriptionDiagnosticSettingsClient.Delete`

```go
ctx := context.TODO()
id := subscriptiondiagnosticsettings.NewDiagnosticSettingID("12345678-1234-9876-4563-123456789012", "diagnosticSettingName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionDiagnosticSettingsClient.Get`

```go
ctx := context.TODO()
id := subscriptiondiagnosticsettings.NewDiagnosticSettingID("12345678-1234-9876-4563-123456789012", "diagnosticSettingName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionDiagnosticSettingsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
