
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/fqdnlistlocalrulestack` Documentation

The `fqdnlistlocalrulestack` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2022-08-29`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/fqdnlistlocalrulestack"
```


### Client Initialization

```go
client := fqdnlistlocalrulestack.NewFqdnListLocalRulestackClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FqdnListLocalRulestackClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := fqdnlistlocalrulestack.NewLocalRulestackFqdnListID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "fqdnListName")

payload := fqdnlistlocalrulestack.FqdnListLocalRulestackResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FqdnListLocalRulestackClient.Delete`

```go
ctx := context.TODO()
id := fqdnlistlocalrulestack.NewLocalRulestackFqdnListID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "fqdnListName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FqdnListLocalRulestackClient.Get`

```go
ctx := context.TODO()
id := fqdnlistlocalrulestack.NewLocalRulestackFqdnListID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "fqdnListName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FqdnListLocalRulestackClient.ListByLocalRulestacks`

```go
ctx := context.TODO()
id := fqdnlistlocalrulestack.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.ListByLocalRulestacks(ctx, id)` can be used to do batched pagination
items, err := client.ListByLocalRulestacksComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
