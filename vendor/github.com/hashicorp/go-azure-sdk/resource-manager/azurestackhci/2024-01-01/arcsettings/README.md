
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/arcsettings` Documentation

The `arcsettings` SDK allows for interaction with Azure Resource Manager `azurestackhci` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/arcsettings"
```


### Client Initialization

```go
client := arcsettings.NewArcSettingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ArcSettingsClient.ArcSettingsCreate`

```go
ctx := context.TODO()
id := arcsettings.NewArcSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "arcSettingName")

payload := arcsettings.ArcSetting{
	// ...
}


read, err := client.ArcSettingsCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ArcSettingsClient.ArcSettingsDelete`

```go
ctx := context.TODO()
id := arcsettings.NewArcSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "arcSettingName")

if err := client.ArcSettingsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ArcSettingsClient.ArcSettingsGet`

```go
ctx := context.TODO()
id := arcsettings.NewArcSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "arcSettingName")

read, err := client.ArcSettingsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ArcSettingsClient.ArcSettingsListByCluster`

```go
ctx := context.TODO()
id := arcsettings.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

// alternatively `client.ArcSettingsListByCluster(ctx, id)` can be used to do batched pagination
items, err := client.ArcSettingsListByClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ArcSettingsClient.ArcSettingsUpdate`

```go
ctx := context.TODO()
id := arcsettings.NewArcSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "arcSettingName")

payload := arcsettings.ArcSettingsPatch{
	// ...
}


read, err := client.ArcSettingsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ArcSettingsClient.ConsentAndInstallDefaultExtensions`

```go
ctx := context.TODO()
id := arcsettings.NewArcSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "arcSettingName")

read, err := client.ConsentAndInstallDefaultExtensions(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ArcSettingsClient.CreateIdentity`

```go
ctx := context.TODO()
id := arcsettings.NewArcSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "arcSettingName")

if err := client.CreateIdentityThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ArcSettingsClient.GeneratePassword`

```go
ctx := context.TODO()
id := arcsettings.NewArcSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "arcSettingName")

read, err := client.GeneratePassword(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ArcSettingsClient.InitializeDisableProcess`

```go
ctx := context.TODO()
id := arcsettings.NewArcSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "arcSettingName")

if err := client.InitializeDisableProcessThenPoll(ctx, id); err != nil {
	// handle the error
}
```
