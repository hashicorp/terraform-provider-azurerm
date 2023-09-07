
## `github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/availabilitygrouplisteners` Documentation

The `availabilitygrouplisteners` SDK allows for interaction with the Azure Resource Manager Service `sqlvirtualmachine` (API Version `2022-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/availabilitygrouplisteners"
```


### Client Initialization

```go
client := availabilitygrouplisteners.NewAvailabilityGroupListenersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AvailabilityGroupListenersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := availabilitygrouplisteners.NewAvailabilityGroupListenerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineGroupValue", "availabilityGroupListenerValue")

payload := availabilitygrouplisteners.AvailabilityGroupListener{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AvailabilityGroupListenersClient.Delete`

```go
ctx := context.TODO()
id := availabilitygrouplisteners.NewAvailabilityGroupListenerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineGroupValue", "availabilityGroupListenerValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AvailabilityGroupListenersClient.Get`

```go
ctx := context.TODO()
id := availabilitygrouplisteners.NewAvailabilityGroupListenerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineGroupValue", "availabilityGroupListenerValue")

read, err := client.Get(ctx, id, availabilitygrouplisteners.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AvailabilityGroupListenersClient.ListByGroup`

```go
ctx := context.TODO()
id := availabilitygrouplisteners.NewSqlVirtualMachineGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineGroupValue")

// alternatively `client.ListByGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
