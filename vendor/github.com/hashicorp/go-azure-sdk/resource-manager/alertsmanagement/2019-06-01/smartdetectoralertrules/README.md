
## `github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-06-01/smartdetectoralertrules` Documentation

The `smartdetectoralertrules` SDK allows for interaction with Azure Resource Manager `alertsmanagement` (API Version `2019-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-06-01/smartdetectoralertrules"
```


### Client Initialization

```go
client := smartdetectoralertrules.NewSmartDetectorAlertRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SmartDetectorAlertRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := smartdetectoralertrules.NewSmartDetectorAlertRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "smartDetectorAlertRuleName")

payload := smartdetectoralertrules.AlertRule{
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


### Example Usage: `SmartDetectorAlertRulesClient.Delete`

```go
ctx := context.TODO()
id := smartdetectoralertrules.NewSmartDetectorAlertRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "smartDetectorAlertRuleName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SmartDetectorAlertRulesClient.Get`

```go
ctx := context.TODO()
id := smartdetectoralertrules.NewSmartDetectorAlertRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "smartDetectorAlertRuleName")

read, err := client.Get(ctx, id, smartdetectoralertrules.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SmartDetectorAlertRulesClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id, smartdetectoralertrules.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, smartdetectoralertrules.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SmartDetectorAlertRulesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, smartdetectoralertrules.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, smartdetectoralertrules.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SmartDetectorAlertRulesClient.Patch`

```go
ctx := context.TODO()
id := smartdetectoralertrules.NewSmartDetectorAlertRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "smartDetectorAlertRuleName")

payload := smartdetectoralertrules.AlertRulePatchObject{
	// ...
}


read, err := client.Patch(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
