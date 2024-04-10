
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/prerules` Documentation

The `prerules` SDK allows for interaction with the Azure Resource Manager Service `paloaltonetworks` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/prerules"
```


### Client Initialization

```go
client := prerules.NewPreRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PreRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := prerules.NewPreRuleID("globalRulestackValue", "preRuleValue")

payload := prerules.PreRulesResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PreRulesClient.Delete`

```go
ctx := context.TODO()
id := prerules.NewPreRuleID("globalRulestackValue", "preRuleValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PreRulesClient.Get`

```go
ctx := context.TODO()
id := prerules.NewPreRuleID("globalRulestackValue", "preRuleValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PreRulesClient.GetCounters`

```go
ctx := context.TODO()
id := prerules.NewPreRuleID("globalRulestackValue", "preRuleValue")

read, err := client.GetCounters(ctx, id, prerules.DefaultGetCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PreRulesClient.List`

```go
ctx := context.TODO()
id := prerules.NewGlobalRulestackID("globalRulestackValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PreRulesClient.RefreshCounters`

```go
ctx := context.TODO()
id := prerules.NewPreRuleID("globalRulestackValue", "preRuleValue")

read, err := client.RefreshCounters(ctx, id, prerules.DefaultRefreshCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PreRulesClient.ResetCounters`

```go
ctx := context.TODO()
id := prerules.NewPreRuleID("globalRulestackValue", "preRuleValue")

read, err := client.ResetCounters(ctx, id, prerules.DefaultResetCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
