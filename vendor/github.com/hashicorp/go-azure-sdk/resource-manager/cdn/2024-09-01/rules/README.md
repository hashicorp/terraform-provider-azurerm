
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-09-01/rules` Documentation

The `rules` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2024-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-09-01/rules"
```


### Client Initialization

```go
client := rules.NewRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RulesClient.Create`

```go
ctx := context.TODO()
id := rules.NewRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "ruleSetName", "ruleName")

payload := rules.Rule{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RulesClient.Delete`

```go
ctx := context.TODO()
id := rules.NewRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "ruleSetName", "ruleName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RulesClient.Get`

```go
ctx := context.TODO()
id := rules.NewRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "ruleSetName", "ruleName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RulesClient.ListByRuleSet`

```go
ctx := context.TODO()
id := rules.NewRuleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "ruleSetName")

// alternatively `client.ListByRuleSet(ctx, id)` can be used to do batched pagination
items, err := client.ListByRuleSetComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RulesClient.Update`

```go
ctx := context.TODO()
id := rules.NewRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "ruleSetName", "ruleName")

payload := rules.RuleUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
