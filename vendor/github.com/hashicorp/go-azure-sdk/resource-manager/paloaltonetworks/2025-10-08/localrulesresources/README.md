
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulesresources` Documentation

The `localrulesresources` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2025-10-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulesresources"
```


### Client Initialization

```go
client := localrulesresources.NewLocalRulesResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LocalRulesResourcesClient.LocalRulesCreateOrUpdate`

```go
ctx := context.TODO()
id := localrulesresources.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "localRuleName")

payload := localrulesresources.LocalRulesResource{
	// ...
}


if err := client.LocalRulesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LocalRulesResourcesClient.LocalRulesDelete`

```go
ctx := context.TODO()
id := localrulesresources.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "localRuleName")

if err := client.LocalRulesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LocalRulesResourcesClient.LocalRulesGet`

```go
ctx := context.TODO()
id := localrulesresources.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "localRuleName")

read, err := client.LocalRulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalRulesResourcesClient.LocalRulesListByLocalRulestacks`

```go
ctx := context.TODO()
id := localrulesresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.LocalRulesListByLocalRulestacks(ctx, id)` can be used to do batched pagination
items, err := client.LocalRulesListByLocalRulestacksComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulesResourcesClient.LocalRulesgetCounters`

```go
ctx := context.TODO()
id := localrulesresources.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "localRuleName")

read, err := client.LocalRulesgetCounters(ctx, id, localrulesresources.DefaultLocalRulesgetCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalRulesResourcesClient.LocalRulesrefreshCounters`

```go
ctx := context.TODO()
id := localrulesresources.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "localRuleName")

read, err := client.LocalRulesrefreshCounters(ctx, id, localrulesresources.DefaultLocalRulesrefreshCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalRulesResourcesClient.LocalRulesresetCounters`

```go
ctx := context.TODO()
id := localrulesresources.NewLocalRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "localRuleName")

read, err := client.LocalRulesresetCounters(ctx, id, localrulesresources.DefaultLocalRulesresetCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
