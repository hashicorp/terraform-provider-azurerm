
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2021-05-01/managedcluster` Documentation

The `managedcluster` SDK allows for interaction with the Azure Resource Manager Service `servicefabricmanagedcluster` (API Version `2021-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2021-05-01/managedcluster"
```


### Client Initialization

```go
client := managedcluster.NewManagedClusterClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedClusterClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := managedcluster.NewManagedClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := managedcluster.ManagedCluster{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedClusterClient.Delete`

```go
ctx := context.TODO()
id := managedcluster.NewManagedClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedClusterClient.Get`

```go
ctx := context.TODO()
id := managedcluster.NewManagedClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedClusterClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := managedcluster.NewResourceGroupID()

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedClusterClient.ListBySubscription`

```go
ctx := context.TODO()
id := managedcluster.NewSubscriptionID()

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedClusterClient.Update`

```go
ctx := context.TODO()
id := managedcluster.NewManagedClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := managedcluster.ManagedClusterUpdateParameters{
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
