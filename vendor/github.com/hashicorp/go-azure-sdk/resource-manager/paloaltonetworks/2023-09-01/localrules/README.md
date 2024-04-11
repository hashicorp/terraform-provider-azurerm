
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/localrules` Documentation

The `localrules` SDK allows for interaction with the Azure Resource Manager Service `paloaltonetworks` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/localrules"
```


### Client Initialization

```go
client := localrules.NewLocalRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LocalRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := localrules.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackValue", "localRuleValue")

payload := localrules.LocalRulesResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LocalRulesClient.Delete`

```go
ctx := context.TODO()
id := localrules.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackValue", "localRuleValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LocalRulesClient.Get`

```go
ctx := context.TODO()
id := localrules.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackValue", "localRuleValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalRulesClient.GetCounters`

```go
ctx := context.TODO()
id := localrules.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackValue", "localRuleValue")

read, err := client.GetCounters(ctx, id, localrules.DefaultGetCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalRulesClient.ListByLocalRulestacks`

```go
ctx := context.TODO()
id := localrules.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackValue")

// alternatively `client.ListByLocalRulestacks(ctx, id)` can be used to do batched pagination
items, err := client.ListByLocalRulestacksComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulesClient.RefreshCounters`

```go
ctx := context.TODO()
id := localrules.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackValue", "localRuleValue")

read, err := client.RefreshCounters(ctx, id, localrules.DefaultRefreshCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalRulesClient.ResetCounters`

```go
ctx := context.TODO()
id := localrules.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackValue", "localRuleValue")

read, err := client.ResetCounters(ctx, id, localrules.DefaultResetCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
