
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/virtualmachineinstances` Documentation

The `virtualmachineinstances` SDK allows for interaction with Azure Resource Manager `azurestackhci` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/virtualmachineinstances"
```


### Client Initialization

```go
client := virtualmachineinstances.NewVirtualMachineInstancesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualMachineInstancesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := virtualmachineinstances.VirtualMachineInstance{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineInstancesClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineInstancesClient.Get`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineInstancesClient.List`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineInstancesClient.Restart`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

if err := client.RestartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineInstancesClient.Start`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineInstancesClient.Stop`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

if err := client.StopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineInstancesClient.Update`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := virtualmachineinstances.VirtualMachineInstanceUpdateRequest{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
