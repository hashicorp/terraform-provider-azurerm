
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups` Documentation

The `capacityreservationgroups` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2022-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
```


### Client Initialization

```go
client := capacityreservationgroups.NewCapacityReservationGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CapacityReservationGroupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := capacityreservationgroups.NewCapacityReservationGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityReservationGroupName")

payload := capacityreservationgroups.CapacityReservationGroup{
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


### Example Usage: `CapacityReservationGroupsClient.Delete`

```go
ctx := context.TODO()
id := capacityreservationgroups.NewCapacityReservationGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityReservationGroupName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CapacityReservationGroupsClient.Get`

```go
ctx := context.TODO()
id := capacityreservationgroups.NewCapacityReservationGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityReservationGroupName")

read, err := client.Get(ctx, id, capacityreservationgroups.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CapacityReservationGroupsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, capacityreservationgroups.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, capacityreservationgroups.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CapacityReservationGroupsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, capacityreservationgroups.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, capacityreservationgroups.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CapacityReservationGroupsClient.Update`

```go
ctx := context.TODO()
id := capacityreservationgroups.NewCapacityReservationGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityReservationGroupName")

payload := capacityreservationgroups.CapacityReservationGroupUpdate{
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
