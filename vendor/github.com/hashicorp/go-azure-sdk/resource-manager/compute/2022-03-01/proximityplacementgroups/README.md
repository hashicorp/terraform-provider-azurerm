
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups` Documentation

The `proximityplacementgroups` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2022-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
```


### Client Initialization

```go
client := proximityplacementgroups.NewProximityPlacementGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProximityPlacementGroupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := proximityplacementgroups.NewProximityPlacementGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "proximityPlacementGroupName")

payload := proximityplacementgroups.ProximityPlacementGroup{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProximityPlacementGroupsClient.Delete`

```go
ctx := context.TODO()
id := proximityplacementgroups.NewProximityPlacementGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "proximityPlacementGroupName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProximityPlacementGroupsClient.Get`

```go
ctx := context.TODO()
id := proximityplacementgroups.NewProximityPlacementGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "proximityPlacementGroupName")

read, err := client.Get(ctx, id, proximityplacementgroups.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProximityPlacementGroupsClient.ListByResourceGroup`

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


### Example Usage: `ProximityPlacementGroupsClient.ListBySubscription`

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


### Example Usage: `ProximityPlacementGroupsClient.Update`

```go
ctx := context.TODO()
id := proximityplacementgroups.NewProximityPlacementGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "proximityPlacementGroupName")

payload := proximityplacementgroups.UpdateResource{
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
