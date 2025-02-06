
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/extensions` Documentation

The `extensions` SDK allows for interaction with Azure Resource Manager `azurestackhci` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/extensions"
```


### Client Initialization

```go
client := extensions.NewExtensionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExtensionsClient.Create`

```go
ctx := context.TODO()
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "arcSettingName", "extensionName")

payload := extensions.Extension{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.Delete`

```go
ctx := context.TODO()
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "arcSettingName", "extensionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.Get`

```go
ctx := context.TODO()
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "arcSettingName", "extensionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExtensionsClient.ListByArcSetting`

```go
ctx := context.TODO()
id := extensions.NewArcSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "arcSettingName")

// alternatively `client.ListByArcSetting(ctx, id)` can be used to do batched pagination
items, err := client.ListByArcSettingComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ExtensionsClient.Update`

```go
ctx := context.TODO()
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "arcSettingName", "extensionName")

payload := extensions.ExtensionPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.Upgrade`

```go
ctx := context.TODO()
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "arcSettingName", "extensionName")

payload := extensions.ExtensionUpgradeParameters{
	// ...
}


if err := client.UpgradeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
