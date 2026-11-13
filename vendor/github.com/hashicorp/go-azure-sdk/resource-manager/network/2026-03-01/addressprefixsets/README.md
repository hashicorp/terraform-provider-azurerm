
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/addressprefixsets` Documentation

The `addressprefixsets` SDK allows for interaction with Azure Resource Manager `network` (API Version `2026-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/addressprefixsets"
```


### Client Initialization

```go
client := addressprefixsets.NewAddressPrefixSetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AddressPrefixSetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := addressprefixsets.NewAddressPrefixSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationSecurityGroupName", "addressPrefixSetName")

payload := addressprefixsets.AddressPrefixSet{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AddressPrefixSetsClient.Delete`

```go
ctx := context.TODO()
id := addressprefixsets.NewAddressPrefixSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationSecurityGroupName", "addressPrefixSetName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AddressPrefixSetsClient.Get`

```go
ctx := context.TODO()
id := addressprefixsets.NewAddressPrefixSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationSecurityGroupName", "addressPrefixSetName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AddressPrefixSetsClient.List`

```go
ctx := context.TODO()
id := addressprefixsets.NewApplicationSecurityGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationSecurityGroupName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
