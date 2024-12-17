
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets` Documentation

The `rulesets` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2024-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets"
```


### Client Initialization

```go
client := rulesets.NewRuleSetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RuleSetsClient.Create`

```go
ctx := context.TODO()
id := rulesets.NewRuleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "ruleSetName")

read, err := client.Create(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RuleSetsClient.Delete`

```go
ctx := context.TODO()
id := rulesets.NewRuleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "ruleSetName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RuleSetsClient.Get`

```go
ctx := context.TODO()
id := rulesets.NewRuleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "ruleSetName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RuleSetsClient.ListByProfile`

```go
ctx := context.TODO()
id := rulesets.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

// alternatively `client.ListByProfile(ctx, id)` can be used to do batched pagination
items, err := client.ListByProfileComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RuleSetsClient.ListResourceUsage`

```go
ctx := context.TODO()
id := rulesets.NewRuleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "ruleSetName")

// alternatively `client.ListResourceUsage(ctx, id)` can be used to do batched pagination
items, err := client.ListResourceUsageComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
