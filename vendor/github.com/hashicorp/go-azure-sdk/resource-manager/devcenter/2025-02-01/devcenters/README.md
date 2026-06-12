
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devcenters` Documentation

The `devcenters` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devcenters"
```


### Client Initialization

```go
client := devcenters.NewDevCentersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DevCentersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := devcenters.NewDevCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName")

payload := devcenters.DevCenter{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DevCentersClient.Delete`

```go
ctx := context.TODO()
id := devcenters.NewDevCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DevCentersClient.Get`

```go
ctx := context.TODO()
id := devcenters.NewDevCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DevCentersClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, devcenters.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, devcenters.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DevCentersClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, devcenters.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, devcenters.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DevCentersClient.Update`

```go
ctx := context.TODO()
id := devcenters.NewDevCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName")

payload := devcenters.DevCenterUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
