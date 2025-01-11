
## `github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/spacecraft` Documentation

The `spacecraft` SDK allows for interaction with Azure Resource Manager `orbital` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/spacecraft"
```


### Client Initialization

```go
client := spacecraft.NewSpacecraftClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SpacecraftClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := spacecraft.NewSpacecraftID("12345678-1234-9876-4563-123456789012", "example-resource-group", "spacecraftName")

payload := spacecraft.Spacecraft{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SpacecraftClient.Delete`

```go
ctx := context.TODO()
id := spacecraft.NewSpacecraftID("12345678-1234-9876-4563-123456789012", "example-resource-group", "spacecraftName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SpacecraftClient.Get`

```go
ctx := context.TODO()
id := spacecraft.NewSpacecraftID("12345678-1234-9876-4563-123456789012", "example-resource-group", "spacecraftName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SpacecraftClient.List`

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


### Example Usage: `SpacecraftClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SpacecraftClient.UpdateTags`

```go
ctx := context.TODO()
id := spacecraft.NewSpacecraftID("12345678-1234-9876-4563-123456789012", "example-resource-group", "spacecraftName")

payload := spacecraft.TagsObject{
	// ...
}


if err := client.UpdateTagsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
