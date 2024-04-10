
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/fqdnlistglobalrulestack` Documentation

The `fqdnlistglobalrulestack` SDK allows for interaction with the Azure Resource Manager Service `paloaltonetworks` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/fqdnlistglobalrulestack"
```


### Client Initialization

```go
client := fqdnlistglobalrulestack.NewFqdnListGlobalRulestackClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FqdnListGlobalRulestackClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := fqdnlistglobalrulestack.NewFqdnListID("globalRulestackValue", "fqdnListValue")

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
id := fqdnlistglobalrulestack.NewFqdnListID("globalRulestackValue", "fqdnListValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FqdnListGlobalRulestackClient.Get`

```go
ctx := context.TODO()
id := fqdnlistglobalrulestack.NewFqdnListID("globalRulestackValue", "fqdnListValue")

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
id := fqdnlistglobalrulestack.NewGlobalRulestackID("globalRulestackValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
