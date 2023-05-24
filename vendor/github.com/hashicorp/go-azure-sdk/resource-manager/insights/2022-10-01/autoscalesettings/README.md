
## `github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-10-01/autoscalesettings` Documentation

The `autoscalesettings` SDK allows for interaction with the Azure Resource Manager Service `insights` (API Version `2022-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-10-01/autoscalesettings"
```


### Client Initialization

```go
client := autoscalesettings.NewAutoScaleSettingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AutoScaleSettingsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := autoscalesettings.NewAutoScaleSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autoScaleSettingValue")

payload := autoscalesettings.AutoscaleSettingResource{
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


### Example Usage: `AutoScaleSettingsClient.Delete`

```go
ctx := context.TODO()
id := autoscalesettings.NewAutoScaleSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autoScaleSettingValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AutoScaleSettingsClient.Get`

```go
ctx := context.TODO()
id := autoscalesettings.NewAutoScaleSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autoScaleSettingValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AutoScaleSettingsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := autoscalesettings.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AutoScaleSettingsClient.ListBySubscription`

```go
ctx := context.TODO()
id := autoscalesettings.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
