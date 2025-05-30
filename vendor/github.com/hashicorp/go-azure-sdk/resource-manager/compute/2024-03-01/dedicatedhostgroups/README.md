
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/dedicatedhostgroups` Documentation

The `dedicatedhostgroups` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/dedicatedhostgroups"
```


### Client Initialization

```go
client := dedicatedhostgroups.NewDedicatedHostGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DedicatedHostGroupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewDedicatedHostGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostGroupName")

payload := dedicatedhostgroups.DedicatedHostGroup{
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


### Example Usage: `DedicatedHostGroupsClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewDedicatedHostGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostGroupName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DedicatedHostGroupsClient.Get`

```go
ctx := context.TODO()
id := commonids.NewDedicatedHostGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostGroupName")

read, err := client.Get(ctx, id, dedicatedhostgroups.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DedicatedHostGroupsClient.ListByResourceGroup`

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


### Example Usage: `DedicatedHostGroupsClient.ListBySubscription`

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


### Example Usage: `DedicatedHostGroupsClient.Update`

```go
ctx := context.TODO()
id := commonids.NewDedicatedHostGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostGroupName")

payload := dedicatedhostgroups.DedicatedHostGroupUpdate{
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
