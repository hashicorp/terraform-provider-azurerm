
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressroutelags` Documentation

The `expressroutelags` SDK allows for interaction with Azure Resource Manager `network` (API Version `2026-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressroutelags"
```


### Client Initialization

```go
client := expressroutelags.NewExpressRouteLagsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteLagsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := expressroutelags.NewExpressRouteLagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteLagName")

payload := expressroutelags.ExpressRouteLag{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteLagsClient.Delete`

```go
ctx := context.TODO()
id := expressroutelags.NewExpressRouteLagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteLagName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteLagsClient.GenerateLoa`

```go
ctx := context.TODO()
id := expressroutelags.NewExpressRouteLagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteLagName")

payload := expressroutelags.GenerateExpressRouteLagsLOARequest{
	// ...
}


read, err := client.GenerateLoa(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteLagsClient.Get`

```go
ctx := context.TODO()
id := expressroutelags.NewExpressRouteLagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteLagName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteLagsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ExpressRouteLagsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ExpressRouteLagsClient.Update`

```go
ctx := context.TODO()
id := expressroutelags.NewExpressRouteLagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteLagName")

payload := expressroutelags.ExpressRouteLagUpdateTagsOrIdentityRequest{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
