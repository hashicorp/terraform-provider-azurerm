
## `github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances` Documentation

The `appliances` SDK allows for interaction with Azure Resource Manager `resourceconnector` (API Version `2022-10-27`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
```


### Client Initialization

```go
client := appliances.NewAppliancesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AppliancesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := appliances.NewApplianceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applianceName")

payload := appliances.Appliance{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppliancesClient.Delete`

```go
ctx := context.TODO()
id := appliances.NewApplianceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applianceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppliancesClient.Get`

```go
ctx := context.TODO()
id := appliances.NewApplianceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applianceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppliancesClient.GetTelemetryConfig`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.GetTelemetryConfig(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppliancesClient.GetUpgradeGraph`

```go
ctx := context.TODO()
id := appliances.NewUpgradeGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applianceName", "upgradeGraphName")

read, err := client.GetUpgradeGraph(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppliancesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppliancesClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppliancesClient.ListClusterUserCredential`

```go
ctx := context.TODO()
id := appliances.NewApplianceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applianceName")

read, err := client.ListClusterUserCredential(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppliancesClient.ListKeys`

```go
ctx := context.TODO()
id := appliances.NewApplianceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applianceName")

read, err := client.ListKeys(ctx, id, appliances.DefaultListKeysOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppliancesClient.Update`

```go
ctx := context.TODO()
id := appliances.NewApplianceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applianceName")

payload := appliances.PatchableAppliance{
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
