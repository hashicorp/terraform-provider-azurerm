
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/prefixlistlocalrulestack` Documentation

The `prefixlistlocalrulestack` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/prefixlistlocalrulestack"
```


### Client Initialization

```go
client := prefixlistlocalrulestack.NewPrefixListLocalRulestackClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrefixListLocalRulestackClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := prefixlistlocalrulestack.NewLocalRulestackPrefixListID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "prefixListName")

payload := prefixlistlocalrulestack.PrefixListResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PrefixListLocalRulestackClient.Delete`

```go
ctx := context.TODO()
id := prefixlistlocalrulestack.NewLocalRulestackPrefixListID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "prefixListName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PrefixListLocalRulestackClient.Get`

```go
ctx := context.TODO()
id := prefixlistlocalrulestack.NewLocalRulestackPrefixListID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "prefixListName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrefixListLocalRulestackClient.ListByLocalRulestacks`

```go
ctx := context.TODO()
id := prefixlistlocalrulestack.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.ListByLocalRulestacks(ctx, id)` can be used to do batched pagination
items, err := client.ListByLocalRulestacksComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
