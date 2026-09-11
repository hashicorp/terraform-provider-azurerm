
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins` Documentation

The `afdorigins` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2025-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins"
```


### Client Initialization

```go
client := afdorigins.NewAFDOriginsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AFDOriginsClient.Create`

```go
ctx := context.TODO()
id := afdorigins.NewOriginGroupOriginID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "originGroupName", "originName")

payload := afdorigins.AFDOrigin{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AFDOriginsClient.Delete`

```go
ctx := context.TODO()
id := afdorigins.NewOriginGroupOriginID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "originGroupName", "originName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AFDOriginsClient.Get`

```go
ctx := context.TODO()
id := afdorigins.NewOriginGroupOriginID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "originGroupName", "originName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AFDOriginsClient.ListByOriginGroup`

```go
ctx := context.TODO()
id := afdorigins.NewOriginGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "originGroupName")

// alternatively `client.ListByOriginGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByOriginGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AFDOriginsClient.Update`

```go
ctx := context.TODO()
id := afdorigins.NewOriginGroupOriginID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "originGroupName", "originName")

payload := afdorigins.AFDOriginUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
