
## `github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2021-08-08/alertsmanagement` Documentation

The `alertsmanagement` SDK allows for interaction with the Azure Resource Manager Service `alertsmanagement` (API Version `2021-08-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2021-08-08/alertsmanagement"
```


### Client Initialization

```go
client := alertsmanagement.NewAlertsManagementClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AlertsManagementClient.AlertProcessingRulesCreateOrUpdate`

```go
ctx := context.TODO()
id := alertsmanagement.NewActionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "alertProcessingRuleValue")

payload := alertsmanagement.AlertProcessingRule{
	// ...
}


read, err := client.AlertProcessingRulesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertsManagementClient.AlertProcessingRulesDelete`

```go
ctx := context.TODO()
id := alertsmanagement.NewActionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "alertProcessingRuleValue")

read, err := client.AlertProcessingRulesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertsManagementClient.AlertProcessingRulesGetByName`

```go
ctx := context.TODO()
id := alertsmanagement.NewActionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "alertProcessingRuleValue")

read, err := client.AlertProcessingRulesGetByName(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertsManagementClient.AlertProcessingRulesListByResourceGroup`

```go
ctx := context.TODO()
id := alertsmanagement.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.AlertProcessingRulesListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.AlertProcessingRulesListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AlertsManagementClient.AlertProcessingRulesListBySubscription`

```go
ctx := context.TODO()
id := alertsmanagement.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.AlertProcessingRulesListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.AlertProcessingRulesListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AlertsManagementClient.AlertProcessingRulesUpdate`

```go
ctx := context.TODO()
id := alertsmanagement.NewActionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "alertProcessingRuleValue")

payload := alertsmanagement.PatchObject{
	// ...
}


read, err := client.AlertProcessingRulesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
