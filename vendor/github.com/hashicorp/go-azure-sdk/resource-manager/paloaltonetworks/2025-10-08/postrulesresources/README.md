
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/postrulesresources` Documentation

The `postrulesresources` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2025-10-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/postrulesresources"
```


### Client Initialization

```go
client := postrulesresources.NewPostRulesResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PostRulesResourcesClient.PostRulesCreateOrUpdate`

```go
ctx := context.TODO()
id := postrulesresources.NewPostRuleID("globalRulestackName", "postRuleName")

payload := postrulesresources.PostRulesResource{
	// ...
}


if err := client.PostRulesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PostRulesResourcesClient.PostRulesDelete`

```go
ctx := context.TODO()
id := postrulesresources.NewPostRuleID("globalRulestackName", "postRuleName")

if err := client.PostRulesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PostRulesResourcesClient.PostRulesGet`

```go
ctx := context.TODO()
id := postrulesresources.NewPostRuleID("globalRulestackName", "postRuleName")

read, err := client.PostRulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PostRulesResourcesClient.PostRulesList`

```go
ctx := context.TODO()
id := postrulesresources.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.PostRulesList(ctx, id)` can be used to do batched pagination
items, err := client.PostRulesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PostRulesResourcesClient.PostRulesgetCounters`

```go
ctx := context.TODO()
id := postrulesresources.NewPostRuleID("globalRulestackName", "postRuleName")

read, err := client.PostRulesgetCounters(ctx, id, postrulesresources.DefaultPostRulesgetCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PostRulesResourcesClient.PostRulesrefreshCounters`

```go
ctx := context.TODO()
id := postrulesresources.NewPostRuleID("globalRulestackName", "postRuleName")

read, err := client.PostRulesrefreshCounters(ctx, id, postrulesresources.DefaultPostRulesrefreshCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PostRulesResourcesClient.PostRulesresetCounters`

```go
ctx := context.TODO()
id := postrulesresources.NewPostRuleID("globalRulestackName", "postRuleName")

read, err := client.PostRulesresetCounters(ctx, id, postrulesresources.DefaultPostRulesresetCountersOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
