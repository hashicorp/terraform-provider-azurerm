
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/customipprefixes` Documentation

The `customipprefixes` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/customipprefixes"
```


### Client Initialization

```go
client := customipprefixes.NewCustomIPPrefixesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CustomIPPrefixesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := customipprefixes.NewCustomIPPrefixID("12345678-1234-9876-4563-123456789012", "example-resource-group", "customIPPrefixValue")

payload := customipprefixes.CustomIPPrefix{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CustomIPPrefixesClient.Delete`

```go
ctx := context.TODO()
id := customipprefixes.NewCustomIPPrefixID("12345678-1234-9876-4563-123456789012", "example-resource-group", "customIPPrefixValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CustomIPPrefixesClient.Get`

```go
ctx := context.TODO()
id := customipprefixes.NewCustomIPPrefixID("12345678-1234-9876-4563-123456789012", "example-resource-group", "customIPPrefixValue")

read, err := client.Get(ctx, id, customipprefixes.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CustomIPPrefixesClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CustomIPPrefixesClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CustomIPPrefixesClient.UpdateTags`

```go
ctx := context.TODO()
id := customipprefixes.NewCustomIPPrefixID("12345678-1234-9876-4563-123456789012", "example-resource-group", "customIPPrefixValue")

payload := customipprefixes.TagsObject{
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
