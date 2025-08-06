
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/virtualmachineruncommands` Documentation

The `virtualmachineruncommands` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/virtualmachineruncommands"
```


### Client Initialization

```go
client := virtualmachineruncommands.NewVirtualMachineRunCommandsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualMachineRunCommandsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualmachineruncommands.NewVirtualMachineRunCommandID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "runCommandName")

payload := virtualmachineruncommands.VirtualMachineRunCommand{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineRunCommandsClient.Delete`

```go
ctx := context.TODO()
id := virtualmachineruncommands.NewVirtualMachineRunCommandID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "runCommandName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineRunCommandsClient.Get`

```go
ctx := context.TODO()
id := virtualmachineruncommands.NewRunCommandID("12345678-1234-9876-4563-123456789012", "locationName", "commandId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineRunCommandsClient.GetByVirtualMachine`

```go
ctx := context.TODO()
id := virtualmachineruncommands.NewVirtualMachineRunCommandID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "runCommandName")

read, err := client.GetByVirtualMachine(ctx, id, virtualmachineruncommands.DefaultGetByVirtualMachineOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineRunCommandsClient.List`

```go
ctx := context.TODO()
id := virtualmachineruncommands.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineRunCommandsClient.ListByVirtualMachine`

```go
ctx := context.TODO()
id := virtualmachineruncommands.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

// alternatively `client.ListByVirtualMachine(ctx, id, virtualmachineruncommands.DefaultListByVirtualMachineOperationOptions())` can be used to do batched pagination
items, err := client.ListByVirtualMachineComplete(ctx, id, virtualmachineruncommands.DefaultListByVirtualMachineOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineRunCommandsClient.Update`

```go
ctx := context.TODO()
id := virtualmachineruncommands.NewVirtualMachineRunCommandID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "runCommandName")

payload := virtualmachineruncommands.VirtualMachineRunCommandUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
