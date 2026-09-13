
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/fqdnlistglobalrulestackresources` Documentation

The `fqdnlistglobalrulestackresources` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2025-10-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/fqdnlistglobalrulestackresources"
```


### Client Initialization

```go
client := fqdnlistglobalrulestackresources.NewFqdnListGlobalRulestackResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FqdnListGlobalRulestackResourcesClient.FqdnListGlobalRulestackCreateOrUpdate`

```go
ctx := context.TODO()
id := fqdnlistglobalrulestackresources.NewFqdnListID("globalRulestackName", "fqdnListName")

payload := fqdnlistglobalrulestackresources.FqdnListGlobalRulestackResource{
	// ...
}


if err := client.FqdnListGlobalRulestackCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FqdnListGlobalRulestackResourcesClient.FqdnListGlobalRulestackDelete`

```go
ctx := context.TODO()
id := fqdnlistglobalrulestackresources.NewFqdnListID("globalRulestackName", "fqdnListName")

if err := client.FqdnListGlobalRulestackDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FqdnListGlobalRulestackResourcesClient.FqdnListGlobalRulestackGet`

```go
ctx := context.TODO()
id := fqdnlistglobalrulestackresources.NewFqdnListID("globalRulestackName", "fqdnListName")

read, err := client.FqdnListGlobalRulestackGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FqdnListGlobalRulestackResourcesClient.FqdnListGlobalRulestackList`

```go
ctx := context.TODO()
id := fqdnlistglobalrulestackresources.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.FqdnListGlobalRulestackList(ctx, id)` can be used to do batched pagination
items, err := client.FqdnListGlobalRulestackListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
