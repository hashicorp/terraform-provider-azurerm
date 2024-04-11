
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/postrules` Documentation

The `postrules` SDK allows for interaction with the Azure Resource Manager Service `paloaltonetworks` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/postrules"
```


### Client Initialization

```go
client := postrules.NewPostRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PostRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := postrules.NewPostRuleID("globalRulestackValue", "postRuleValue")

payload := postrules.PostRulesResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PostRulesClient.Delete`

```go
ctx := context.TODO()
id := postrules.NewPostRuleID("globalRulestackValue", "postRuleValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PostRulesClient.Get`

```go
ctx := context.TODO()
id := postrules.NewPostRuleID("globalRulestackValue", "postRuleValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PostRulesClient.GetCounters`

```go
ctx := context.TODO()
id := postrules.NewPostRuleID("globalRulestackValue", "postRuleValue")

read, err := client.GetCounters(ctx, id, postrules.DefaultGetCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PostRulesClient.List`

```go
ctx := context.TODO()
id := postrules.NewGlobalRulestackID("globalRulestackValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PostRulesClient.RefreshCounters`

```go
ctx := context.TODO()
id := postrules.NewPostRuleID("globalRulestackValue", "postRuleValue")

read, err := client.RefreshCounters(ctx, id, postrules.DefaultRefreshCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PostRulesClient.ResetCounters`

```go
ctx := context.TODO()
id := postrules.NewPostRuleID("globalRulestackValue", "postRuleValue")

read, err := client.ResetCounters(ctx, id, postrules.DefaultResetCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
