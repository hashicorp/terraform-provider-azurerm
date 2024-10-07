
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/ipallocations` Documentation

The `ipallocations` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/ipallocations"
```


### Client Initialization

```go
client := ipallocations.NewIPAllocationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IPAllocationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := ipallocations.NewIPAllocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ipAllocationName")

payload := ipallocations.IPAllocation{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `IPAllocationsClient.Delete`

```go
ctx := context.TODO()
id := ipallocations.NewIPAllocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ipAllocationName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IPAllocationsClient.Get`

```go
ctx := context.TODO()
id := ipallocations.NewIPAllocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ipAllocationName")

read, err := client.Get(ctx, id, ipallocations.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IPAllocationsClient.List`

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


### Example Usage: `IPAllocationsClient.ListByResourceGroup`

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


### Example Usage: `IPAllocationsClient.UpdateTags`

```go
ctx := context.TODO()
id := ipallocations.NewIPAllocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ipAllocationName")

payload := ipallocations.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
