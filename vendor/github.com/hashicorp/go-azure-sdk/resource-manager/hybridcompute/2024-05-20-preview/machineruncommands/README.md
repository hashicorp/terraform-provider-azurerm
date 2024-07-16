
## `github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-05-20-preview/machineruncommands` Documentation

The `machineruncommands` SDK allows for interaction with the Azure Resource Manager Service `hybridcompute` (API Version `2024-05-20-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-05-20-preview/machineruncommands"
```


### Client Initialization

```go
client := machineruncommands.NewMachineRunCommandsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MachineRunCommandsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := machineruncommands.NewRunCommandID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineValue", "runCommandValue")

payload := machineruncommands.MachineRunCommand{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `MachineRunCommandsClient.Delete`

```go
ctx := context.TODO()
id := machineruncommands.NewRunCommandID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineValue", "runCommandValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `MachineRunCommandsClient.Get`

```go
ctx := context.TODO()
id := machineruncommands.NewRunCommandID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineValue", "runCommandValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MachineRunCommandsClient.List`

```go
ctx := context.TODO()
id := machineruncommands.NewMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineValue")

// alternatively `client.List(ctx, id, machineruncommands.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, machineruncommands.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MachineRunCommandsClient.Update`

```go
ctx := context.TODO()
id := machineruncommands.NewRunCommandID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineValue", "runCommandValue")

payload := machineruncommands.ResourceUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
