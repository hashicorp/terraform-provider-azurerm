
## `github.com/hashicorp/go-azure-sdk/resource-manager/fabric/2023-11-01/fabriccapacities` Documentation

The `fabriccapacities` SDK allows for interaction with Azure Resource Manager `fabric` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/fabric/2023-11-01/fabriccapacities"
```


### Client Initialization

```go
client := fabriccapacities.NewFabricCapacitiesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FabricCapacitiesClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := fabriccapacities.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := fabriccapacities.CheckNameAvailabilityRequest{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FabricCapacitiesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := fabriccapacities.NewCapacityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityName")

payload := fabriccapacities.FabricCapacity{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FabricCapacitiesClient.Delete`

```go
ctx := context.TODO()
id := fabriccapacities.NewCapacityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FabricCapacitiesClient.Get`

```go
ctx := context.TODO()
id := fabriccapacities.NewCapacityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FabricCapacitiesClient.ListByResourceGroup`

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


### Example Usage: `FabricCapacitiesClient.ListBySubscription`

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


### Example Usage: `FabricCapacitiesClient.ListSkus`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListSkus(ctx, id)` can be used to do batched pagination
items, err := client.ListSkusComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FabricCapacitiesClient.ListSkusForCapacity`

```go
ctx := context.TODO()
id := fabriccapacities.NewCapacityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityName")

// alternatively `client.ListSkusForCapacity(ctx, id)` can be used to do batched pagination
items, err := client.ListSkusForCapacityComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FabricCapacitiesClient.Resume`

```go
ctx := context.TODO()
id := fabriccapacities.NewCapacityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityName")

if err := client.ResumeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FabricCapacitiesClient.Suspend`

```go
ctx := context.TODO()
id := fabriccapacities.NewCapacityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityName")

if err := client.SuspendThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FabricCapacitiesClient.Update`

```go
ctx := context.TODO()
id := fabriccapacities.NewCapacityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityName")

payload := fabriccapacities.FabricCapacityUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
