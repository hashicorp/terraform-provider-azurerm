
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskaccesses` Documentation

The `diskaccesses` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2022-03-02`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskaccesses"
```


### Client Initialization

```go
client := diskaccesses.NewDiskAccessesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DiskAccessesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := diskaccesses.NewDiskAccessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskAccessName")

payload := diskaccesses.DiskAccess{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DiskAccessesClient.Delete`

```go
ctx := context.TODO()
id := diskaccesses.NewDiskAccessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskAccessName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DiskAccessesClient.DeleteAPrivateEndpointConnection`

```go
ctx := context.TODO()
id := diskaccesses.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskAccessName", "privateEndpointConnectionName")

if err := client.DeleteAPrivateEndpointConnectionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DiskAccessesClient.Get`

```go
ctx := context.TODO()
id := diskaccesses.NewDiskAccessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskAccessName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiskAccessesClient.GetAPrivateEndpointConnection`

```go
ctx := context.TODO()
id := diskaccesses.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskAccessName", "privateEndpointConnectionName")

read, err := client.GetAPrivateEndpointConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiskAccessesClient.GetPrivateLinkResources`

```go
ctx := context.TODO()
id := diskaccesses.NewDiskAccessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskAccessName")

read, err := client.GetPrivateLinkResources(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiskAccessesClient.List`

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


### Example Usage: `DiskAccessesClient.ListByResourceGroup`

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


### Example Usage: `DiskAccessesClient.ListPrivateEndpointConnections`

```go
ctx := context.TODO()
id := diskaccesses.NewDiskAccessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskAccessName")

// alternatively `client.ListPrivateEndpointConnections(ctx, id)` can be used to do batched pagination
items, err := client.ListPrivateEndpointConnectionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DiskAccessesClient.Update`

```go
ctx := context.TODO()
id := diskaccesses.NewDiskAccessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskAccessName")

payload := diskaccesses.DiskAccessUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DiskAccessesClient.UpdateAPrivateEndpointConnection`

```go
ctx := context.TODO()
id := diskaccesses.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskAccessName", "privateEndpointConnectionName")

payload := diskaccesses.PrivateEndpointConnection{
	// ...
}


if err := client.UpdateAPrivateEndpointConnectionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
