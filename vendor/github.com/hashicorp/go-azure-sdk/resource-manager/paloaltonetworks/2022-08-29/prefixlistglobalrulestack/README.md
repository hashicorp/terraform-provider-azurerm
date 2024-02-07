
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/prefixlistglobalrulestack` Documentation

The `prefixlistglobalrulestack` SDK allows for interaction with the Azure Resource Manager Service `paloaltonetworks` (API Version `2022-08-29`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/prefixlistglobalrulestack"
```


### Client Initialization

```go
client := prefixlistglobalrulestack.NewPrefixListGlobalRulestackClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrefixListGlobalRulestackClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := prefixlistglobalrulestack.NewPrefixListID("globalRulestackValue", "prefixListValue")

payload := prefixlistglobalrulestack.PrefixListGlobalRulestackResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PrefixListGlobalRulestackClient.Delete`

```go
ctx := context.TODO()
id := prefixlistglobalrulestack.NewPrefixListID("globalRulestackValue", "prefixListValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PrefixListGlobalRulestackClient.Get`

```go
ctx := context.TODO()
id := prefixlistglobalrulestack.NewPrefixListID("globalRulestackValue", "prefixListValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrefixListGlobalRulestackClient.List`

```go
ctx := context.TODO()
id := prefixlistglobalrulestack.NewGlobalRulestackID("globalRulestackValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
