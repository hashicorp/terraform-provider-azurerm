
## `github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/securitymlanalyticssettings` Documentation

The `securitymlanalyticssettings` SDK allows for interaction with the Azure Resource Manager Service `securityinsights` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/securitymlanalyticssettings"
```


### Client Initialization

```go
client := securitymlanalyticssettings.NewSecurityMLAnalyticsSettingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SecurityMLAnalyticsSettingsClient.SecurityMLAnalyticsSettingsCreateOrUpdate`

```go
ctx := context.TODO()
id := securitymlanalyticssettings.NewSecurityMLAnalyticsSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "securityMLAnalyticsSettingValue")

payload := securitymlanalyticssettings.SecurityMLAnalyticsSetting{
	// ...
}


read, err := client.SecurityMLAnalyticsSettingsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecurityMLAnalyticsSettingsClient.SecurityMLAnalyticsSettingsDelete`

```go
ctx := context.TODO()
id := securitymlanalyticssettings.NewSecurityMLAnalyticsSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "securityMLAnalyticsSettingValue")

read, err := client.SecurityMLAnalyticsSettingsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecurityMLAnalyticsSettingsClient.SecurityMLAnalyticsSettingsGet`

```go
ctx := context.TODO()
id := securitymlanalyticssettings.NewSecurityMLAnalyticsSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "securityMLAnalyticsSettingValue")

read, err := client.SecurityMLAnalyticsSettingsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecurityMLAnalyticsSettingsClient.SecurityMLAnalyticsSettingsList`

```go
ctx := context.TODO()
id := securitymlanalyticssettings.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.SecurityMLAnalyticsSettingsList(ctx, id)` can be used to do batched pagination
items, err := client.SecurityMLAnalyticsSettingsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
