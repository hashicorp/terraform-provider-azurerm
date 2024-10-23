
## `github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers` Documentation

The `vmmservers` SDK allows for interaction with Azure Resource Manager `systemcentervirtualmachinemanager` (API Version `2023-10-07`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
```


### Client Initialization

```go
client := vmmservers.NewVMmServersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VMmServersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := vmmservers.NewVMmServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vmmServerName")

payload := vmmservers.VMmServer{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VMmServersClient.Delete`

```go
ctx := context.TODO()
id := vmmservers.NewVMmServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vmmServerName")

if err := client.DeleteThenPoll(ctx, id, vmmservers.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VMmServersClient.Get`

```go
ctx := context.TODO()
id := vmmservers.NewVMmServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vmmServerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VMmServersClient.ListByResourceGroup`

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


### Example Usage: `VMmServersClient.ListBySubscription`

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


### Example Usage: `VMmServersClient.Update`

```go
ctx := context.TODO()
id := vmmservers.NewVMmServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vmmServerName")

payload := vmmservers.VMmServerTagsUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
