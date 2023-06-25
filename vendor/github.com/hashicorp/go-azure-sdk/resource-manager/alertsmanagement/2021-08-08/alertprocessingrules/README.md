
## `github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2021-08-08/alertprocessingrules` Documentation

The `alertprocessingrules` SDK allows for interaction with the Azure Resource Manager Service `alertsmanagement` (API Version `2021-08-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2021-08-08/alertprocessingrules"
```


### Client Initialization

```go
client := alertprocessingrules.NewAlertProcessingRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AlertProcessingRulesClient.AlertProcessingRulesCreateOrUpdate`

```go
ctx := context.TODO()
id := alertprocessingrules.NewActionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "actionRuleValue")

payload := alertprocessingrules.AlertProcessingRule{
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


### Example Usage: `AlertProcessingRulesClient.AlertProcessingRulesDelete`

```go
ctx := context.TODO()
id := alertprocessingrules.NewActionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "actionRuleValue")

read, err := client.AlertProcessingRulesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertProcessingRulesClient.AlertProcessingRulesGetByName`

```go
ctx := context.TODO()
id := alertprocessingrules.NewActionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "actionRuleValue")

read, err := client.AlertProcessingRulesGetByName(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertProcessingRulesClient.AlertProcessingRulesListByResourceGroup`

```go
ctx := context.TODO()
id := alertprocessingrules.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.AlertProcessingRulesListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.AlertProcessingRulesListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AlertProcessingRulesClient.AlertProcessingRulesListBySubscription`

```go
ctx := context.TODO()
id := alertprocessingrules.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.AlertProcessingRulesListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.AlertProcessingRulesListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AlertProcessingRulesClient.AlertProcessingRulesUpdate`

```go
ctx := context.TODO()
id := alertprocessingrules.NewActionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "actionRuleValue")

payload := alertprocessingrules.PatchObject{
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
