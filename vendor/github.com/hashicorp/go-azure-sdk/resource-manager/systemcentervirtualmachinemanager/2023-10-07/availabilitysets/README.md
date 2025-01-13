
## `github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/availabilitysets` Documentation

The `availabilitysets` SDK allows for interaction with Azure Resource Manager `systemcentervirtualmachinemanager` (API Version `2023-10-07`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/availabilitysets"
```


### Client Initialization

```go
client := availabilitysets.NewAvailabilitySetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AvailabilitySetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := availabilitysets.NewAvailabilitySetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "availabilitySetName")

payload := availabilitysets.AvailabilitySet{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AvailabilitySetsClient.Delete`

```go
ctx := context.TODO()
id := availabilitysets.NewAvailabilitySetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "availabilitySetName")

if err := client.DeleteThenPoll(ctx, id, availabilitysets.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `AvailabilitySetsClient.Get`

```go
ctx := context.TODO()
id := availabilitysets.NewAvailabilitySetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "availabilitySetName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AvailabilitySetsClient.ListByResourceGroup`

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


### Example Usage: `AvailabilitySetsClient.ListBySubscription`

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


### Example Usage: `AvailabilitySetsClient.Update`

```go
ctx := context.TODO()
id := availabilitysets.NewAvailabilitySetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "availabilitySetName")

payload := availabilitysets.AvailabilitySetTagsUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
