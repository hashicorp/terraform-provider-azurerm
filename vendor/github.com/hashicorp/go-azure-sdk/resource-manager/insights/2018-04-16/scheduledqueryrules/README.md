
## `github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules` Documentation

The `scheduledqueryrules` SDK allows for interaction with the Azure Resource Manager Service `insights` (API Version `2018-04-16`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
```


### Client Initialization

```go
client := scheduledqueryrules.NewScheduledQueryRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ScheduledQueryRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := scheduledqueryrules.NewScheduledQueryRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "scheduledQueryRuleValue")

payload := scheduledqueryrules.LogSearchRuleResource{
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


### Example Usage: `ScheduledQueryRulesClient.Delete`

```go
ctx := context.TODO()
id := scheduledqueryrules.NewScheduledQueryRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "scheduledQueryRuleValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledQueryRulesClient.Get`

```go
ctx := context.TODO()
id := scheduledqueryrules.NewScheduledQueryRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "scheduledQueryRuleValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledQueryRulesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := scheduledqueryrules.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.ListByResourceGroup(ctx, id, scheduledqueryrules.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledQueryRulesClient.ListBySubscription`

```go
ctx := context.TODO()
id := scheduledqueryrules.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListBySubscription(ctx, id, scheduledqueryrules.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledQueryRulesClient.Update`

```go
ctx := context.TODO()
id := scheduledqueryrules.NewScheduledQueryRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "scheduledQueryRuleValue")

payload := scheduledqueryrules.LogSearchRuleResourcePatch{
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
