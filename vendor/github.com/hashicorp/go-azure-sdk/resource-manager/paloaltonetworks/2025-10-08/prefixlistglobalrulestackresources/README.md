
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/prefixlistglobalrulestackresources` Documentation

The `prefixlistglobalrulestackresources` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2025-10-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/prefixlistglobalrulestackresources"
```


### Client Initialization

```go
client := prefixlistglobalrulestackresources.NewPrefixListGlobalRulestackResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrefixListGlobalRulestackResourcesClient.PrefixListGlobalRulestackCreateOrUpdate`

```go
ctx := context.TODO()
id := prefixlistglobalrulestackresources.NewPrefixListID("globalRulestackName", "prefixListName")

payload := prefixlistglobalrulestackresources.PrefixListGlobalRulestackResource{
	// ...
}


if err := client.PrefixListGlobalRulestackCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PrefixListGlobalRulestackResourcesClient.PrefixListGlobalRulestackDelete`

```go
ctx := context.TODO()
id := prefixlistglobalrulestackresources.NewPrefixListID("globalRulestackName", "prefixListName")

if err := client.PrefixListGlobalRulestackDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PrefixListGlobalRulestackResourcesClient.PrefixListGlobalRulestackGet`

```go
ctx := context.TODO()
id := prefixlistglobalrulestackresources.NewPrefixListID("globalRulestackName", "prefixListName")

read, err := client.PrefixListGlobalRulestackGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrefixListGlobalRulestackResourcesClient.PrefixListGlobalRulestackList`

```go
ctx := context.TODO()
id := prefixlistglobalrulestackresources.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.PrefixListGlobalRulestackList(ctx, id)` can be used to do batched pagination
items, err := client.PrefixListGlobalRulestackListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
