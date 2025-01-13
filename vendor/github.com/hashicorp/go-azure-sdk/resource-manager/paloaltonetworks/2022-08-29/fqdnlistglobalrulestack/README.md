
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/fqdnlistglobalrulestack` Documentation

The `fqdnlistglobalrulestack` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2022-08-29`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/fqdnlistglobalrulestack"
```


### Client Initialization

```go
client := fqdnlistglobalrulestack.NewFqdnListGlobalRulestackClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FqdnListGlobalRulestackClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := fqdnlistglobalrulestack.NewFqdnListID("globalRulestackName", "fqdnListName")

payload := fqdnlistglobalrulestack.FqdnListGlobalRulestackResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FqdnListGlobalRulestackClient.Delete`

```go
ctx := context.TODO()
id := fqdnlistglobalrulestack.NewFqdnListID("globalRulestackName", "fqdnListName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FqdnListGlobalRulestackClient.Get`

```go
ctx := context.TODO()
id := fqdnlistglobalrulestack.NewFqdnListID("globalRulestackName", "fqdnListName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FqdnListGlobalRulestackClient.List`

```go
ctx := context.TODO()
id := fqdnlistglobalrulestack.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
