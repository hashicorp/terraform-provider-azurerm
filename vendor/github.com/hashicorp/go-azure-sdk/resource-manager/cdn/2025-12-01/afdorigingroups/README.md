
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups` Documentation

The `afdorigingroups` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2025-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
```


### Client Initialization

```go
client := afdorigingroups.NewAFDOriginGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AFDOriginGroupsClient.Create`

```go
ctx := context.TODO()
id := afdorigingroups.NewOriginGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "originGroupName")

payload := afdorigingroups.AFDOriginGroup{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AFDOriginGroupsClient.Delete`

```go
ctx := context.TODO()
id := afdorigingroups.NewOriginGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "originGroupName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AFDOriginGroupsClient.Get`

```go
ctx := context.TODO()
id := afdorigingroups.NewOriginGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "originGroupName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AFDOriginGroupsClient.ListByProfile`

```go
ctx := context.TODO()
id := afdorigingroups.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

// alternatively `client.ListByProfile(ctx, id)` can be used to do batched pagination
items, err := client.ListByProfileComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AFDOriginGroupsClient.ListResourceUsage`

```go
ctx := context.TODO()
id := afdorigingroups.NewOriginGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "originGroupName")

// alternatively `client.ListResourceUsage(ctx, id)` can be used to do batched pagination
items, err := client.ListResourceUsageComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AFDOriginGroupsClient.Update`

```go
ctx := context.TODO()
id := afdorigingroups.NewOriginGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "originGroupName")

payload := afdorigingroups.AFDOriginGroupUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
