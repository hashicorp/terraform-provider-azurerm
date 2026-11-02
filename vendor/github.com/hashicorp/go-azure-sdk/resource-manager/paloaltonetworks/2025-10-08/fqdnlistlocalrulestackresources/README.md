
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/fqdnlistlocalrulestackresources` Documentation

The `fqdnlistlocalrulestackresources` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2025-10-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/fqdnlistlocalrulestackresources"
```


### Client Initialization

```go
client := fqdnlistlocalrulestackresources.NewFqdnListLocalRulestackResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FqdnListLocalRulestackResourcesClient.FqdnListLocalRulestackCreateOrUpdate`

```go
ctx := context.TODO()
id := fqdnlistlocalrulestackresources.NewLocalRulestackFqdnListID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "fqdnListName")

payload := fqdnlistlocalrulestackresources.FqdnListLocalRulestackResource{
	// ...
}


if err := client.FqdnListLocalRulestackCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FqdnListLocalRulestackResourcesClient.FqdnListLocalRulestackDelete`

```go
ctx := context.TODO()
id := fqdnlistlocalrulestackresources.NewLocalRulestackFqdnListID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "fqdnListName")

if err := client.FqdnListLocalRulestackDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FqdnListLocalRulestackResourcesClient.FqdnListLocalRulestackGet`

```go
ctx := context.TODO()
id := fqdnlistlocalrulestackresources.NewLocalRulestackFqdnListID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "fqdnListName")

read, err := client.FqdnListLocalRulestackGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FqdnListLocalRulestackResourcesClient.FqdnListLocalRulestackListByLocalRulestacks`

```go
ctx := context.TODO()
id := fqdnlistlocalrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.FqdnListLocalRulestackListByLocalRulestacks(ctx, id)` can be used to do batched pagination
items, err := client.FqdnListLocalRulestackListByLocalRulestacksComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
