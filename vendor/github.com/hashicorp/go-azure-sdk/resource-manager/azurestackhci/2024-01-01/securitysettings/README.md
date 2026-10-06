
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/securitysettings` Documentation

The `securitysettings` SDK allows for interaction with Azure Resource Manager `azurestackhci` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/securitysettings"
```


### Client Initialization

```go
client := securitysettings.NewSecuritySettingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SecuritySettingsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := securitysettings.NewSecuritySettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "securitySettingName")

payload := securitysettings.SecuritySetting{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SecuritySettingsClient.Delete`

```go
ctx := context.TODO()
id := securitysettings.NewSecuritySettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "securitySettingName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SecuritySettingsClient.Get`

```go
ctx := context.TODO()
id := securitysettings.NewSecuritySettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "securitySettingName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecuritySettingsClient.ListByClusters`

```go
ctx := context.TODO()
id := securitysettings.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

// alternatively `client.ListByClusters(ctx, id)` can be used to do batched pagination
items, err := client.ListByClustersComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
