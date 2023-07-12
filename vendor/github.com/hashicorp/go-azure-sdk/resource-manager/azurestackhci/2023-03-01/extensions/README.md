
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/extensions` Documentation

The `extensions` SDK allows for interaction with the Azure Resource Manager Service `azurestackhci` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/extensions"
```


### Client Initialization

```go
client := extensions.NewExtensionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExtensionsClient.ExtensionsCreate`

```go
ctx := context.TODO()
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "arcSettingValue", "extensionValue")

payload := extensions.Extension{
	// ...
}


if err := client.ExtensionsCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.ExtensionsDelete`

```go
ctx := context.TODO()
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "arcSettingValue", "extensionValue")

if err := client.ExtensionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.ExtensionsGet`

```go
ctx := context.TODO()
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "arcSettingValue", "extensionValue")

read, err := client.ExtensionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExtensionsClient.ExtensionsListByArcSetting`

```go
ctx := context.TODO()
id := extensions.NewArcSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "arcSettingValue")

// alternatively `client.ExtensionsListByArcSetting(ctx, id)` can be used to do batched pagination
items, err := client.ExtensionsListByArcSettingComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ExtensionsClient.ExtensionsUpdate`

```go
ctx := context.TODO()
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "arcSettingValue", "extensionValue")

payload := extensions.Extension{
	// ...
}


if err := client.ExtensionsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.ExtensionsUpgrade`

```go
ctx := context.TODO()
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "arcSettingValue", "extensionValue")

payload := extensions.ExtensionUpgradeParameters{
	// ...
}


if err := client.ExtensionsUpgradeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
