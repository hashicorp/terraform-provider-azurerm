
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerappssessionpools` Documentation

The `containerappssessionpools` SDK allows for interaction with Azure Resource Manager `containerapps` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerappssessionpools"
```


### Client Initialization

```go
client := containerappssessionpools.NewContainerAppsSessionPoolsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ContainerAppsSessionPoolsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := containerappssessionpools.NewSessionPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sessionPoolName")

payload := containerappssessionpools.SessionPool{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerAppsSessionPoolsClient.Delete`

```go
ctx := context.TODO()
id := containerappssessionpools.NewSessionPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sessionPoolName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerAppsSessionPoolsClient.Get`

```go
ctx := context.TODO()
id := containerappssessionpools.NewSessionPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sessionPoolName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerAppsSessionPoolsClient.ListByResourceGroup`

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


### Example Usage: `ContainerAppsSessionPoolsClient.ListBySubscription`

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


### Example Usage: `ContainerAppsSessionPoolsClient.Update`

```go
ctx := context.TODO()
id := containerappssessionpools.NewSessionPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sessionPoolName")

payload := containerappssessionpools.SessionPoolUpdatableProperties{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
