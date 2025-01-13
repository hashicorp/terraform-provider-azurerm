
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrules` Documentation

The `localrules` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2022-08-29`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrules"
```


### Client Initialization

```go
client := localrules.NewLocalRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LocalRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := localrules.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "localRuleName")

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
id := localrules.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "localRuleName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LocalRulesClient.Get`

```go
ctx := context.TODO()
id := localrules.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "localRuleName")

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
id := localrules.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "localRuleName")

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
id := localrules.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

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
id := localrules.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "localRuleName")

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
id := localrules.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "localRuleName")

read, err := client.ResetCounters(ctx, id, localrules.DefaultResetCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
