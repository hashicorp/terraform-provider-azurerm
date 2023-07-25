
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/diskpools` Documentation

The `diskpools` SDK allows for interaction with the Azure Resource Manager Service `storagepool` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/diskpools"
```


### Client Initialization

```go
client := diskpools.NewDiskPoolsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DiskPoolsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := diskpools.NewDiskPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue")

payload := diskpools.DiskPoolCreate{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DiskPoolsClient.Deallocate`

```go
ctx := context.TODO()
id := diskpools.NewDiskPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue")

if err := client.DeallocateThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DiskPoolsClient.Delete`

```go
ctx := context.TODO()
id := diskpools.NewDiskPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DiskPoolsClient.Get`

```go
ctx := context.TODO()
id := diskpools.NewDiskPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiskPoolsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := diskpools.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DiskPoolsClient.ListBySubscription`

```go
ctx := context.TODO()
id := diskpools.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DiskPoolsClient.ListOutboundNetworkDependenciesEndpoints`

```go
ctx := context.TODO()
id := diskpools.NewDiskPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue")

// alternatively `client.ListOutboundNetworkDependenciesEndpoints(ctx, id)` can be used to do batched pagination
items, err := client.ListOutboundNetworkDependenciesEndpointsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DiskPoolsClient.Start`

```go
ctx := context.TODO()
id := diskpools.NewDiskPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DiskPoolsClient.Update`

```go
ctx := context.TODO()
id := diskpools.NewDiskPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue")

payload := diskpools.DiskPoolUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DiskPoolsClient.Upgrade`

```go
ctx := context.TODO()
id := diskpools.NewDiskPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue")

if err := client.UpgradeThenPoll(ctx, id); err != nil {
	// handle the error
}
```
