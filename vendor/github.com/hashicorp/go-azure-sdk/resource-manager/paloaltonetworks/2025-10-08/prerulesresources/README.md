
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/prerulesresources` Documentation

The `prerulesresources` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2025-10-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/prerulesresources"
```


### Client Initialization

```go
client := prerulesresources.NewPreRulesResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PreRulesResourcesClient.PreRulesCreateOrUpdate`

```go
ctx := context.TODO()
id := prerulesresources.NewPreRuleID("globalRulestackName", "preRuleName")

payload := prerulesresources.PreRulesResource{
	// ...
}


if err := client.PreRulesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PreRulesResourcesClient.PreRulesDelete`

```go
ctx := context.TODO()
id := prerulesresources.NewPreRuleID("globalRulestackName", "preRuleName")

if err := client.PreRulesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PreRulesResourcesClient.PreRulesGet`

```go
ctx := context.TODO()
id := prerulesresources.NewPreRuleID("globalRulestackName", "preRuleName")

read, err := client.PreRulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PreRulesResourcesClient.PreRulesList`

```go
ctx := context.TODO()
id := prerulesresources.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.PreRulesList(ctx, id)` can be used to do batched pagination
items, err := client.PreRulesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PreRulesResourcesClient.PreRulesgetCounters`

```go
ctx := context.TODO()
id := prerulesresources.NewPreRuleID("globalRulestackName", "preRuleName")

read, err := client.PreRulesgetCounters(ctx, id, prerulesresources.DefaultPreRulesgetCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PreRulesResourcesClient.PreRulesrefreshCounters`

```go
ctx := context.TODO()
id := prerulesresources.NewPreRuleID("globalRulestackName", "preRuleName")

read, err := client.PreRulesrefreshCounters(ctx, id, prerulesresources.DefaultPreRulesrefreshCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PreRulesResourcesClient.PreRulesresetCounters`

```go
ctx := context.TODO()
id := prerulesresources.NewPreRuleID("globalRulestackName", "preRuleName")

read, err := client.PreRulesresetCounters(ctx, id, prerulesresources.DefaultPreRulesresetCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
