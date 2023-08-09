
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/publicipprefixes` Documentation

The `publicipprefixes` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/publicipprefixes"
```


### Client Initialization

```go
client := publicipprefixes.NewPublicIPPrefixesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PublicIPPrefixesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := publicipprefixes.NewPublicIPPrefixID("12345678-1234-9876-4563-123456789012", "example-resource-group", "publicIPPrefixValue")

payload := publicipprefixes.PublicIPPrefix{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PublicIPPrefixesClient.Delete`

```go
ctx := context.TODO()
id := publicipprefixes.NewPublicIPPrefixID("12345678-1234-9876-4563-123456789012", "example-resource-group", "publicIPPrefixValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PublicIPPrefixesClient.Get`

```go
ctx := context.TODO()
id := publicipprefixes.NewPublicIPPrefixID("12345678-1234-9876-4563-123456789012", "example-resource-group", "publicIPPrefixValue")

read, err := client.Get(ctx, id, publicipprefixes.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PublicIPPrefixesClient.List`

```go
ctx := context.TODO()
id := publicipprefixes.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PublicIPPrefixesClient.ListAll`

```go
ctx := context.TODO()
id := publicipprefixes.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PublicIPPrefixesClient.UpdateTags`

```go
ctx := context.TODO()
id := publicipprefixes.NewPublicIPPrefixID("12345678-1234-9876-4563-123456789012", "example-resource-group", "publicIPPrefixValue")

payload := publicipprefixes.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
