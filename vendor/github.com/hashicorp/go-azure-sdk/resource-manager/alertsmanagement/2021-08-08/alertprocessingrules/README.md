
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


### Example Usage: `AlertProcessingRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := alertprocessingrules.NewActionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "actionRuleValue")

payload := alertprocessingrules.AlertProcessingRule{
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


### Example Usage: `AlertProcessingRulesClient.Delete`

```go
ctx := context.TODO()
id := alertprocessingrules.NewActionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "actionRuleValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertProcessingRulesClient.GetByName`

```go
ctx := context.TODO()
id := alertprocessingrules.NewActionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "actionRuleValue")

read, err := client.GetByName(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertProcessingRulesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := alertprocessingrules.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AlertProcessingRulesClient.ListBySubscription`

```go
ctx := context.TODO()
id := alertprocessingrules.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AlertProcessingRulesClient.Update`

```go
ctx := context.TODO()
id := alertprocessingrules.NewActionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "actionRuleValue")

payload := alertprocessingrules.PatchObject{
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
