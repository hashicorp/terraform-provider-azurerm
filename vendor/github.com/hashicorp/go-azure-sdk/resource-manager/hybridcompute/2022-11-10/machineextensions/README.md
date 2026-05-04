
## `github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machineextensions` Documentation

The `machineextensions` SDK allows for interaction with Azure Resource Manager `hybridcompute` (API Version `2022-11-10`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machineextensions"
```


### Client Initialization

```go
client := machineextensions.NewMachineExtensionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MachineExtensionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := machineextensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName", "extensionName")

payload := machineextensions.MachineExtension{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `MachineExtensionsClient.Delete`

```go
ctx := context.TODO()
id := machineextensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName", "extensionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `MachineExtensionsClient.Get`

```go
ctx := context.TODO()
id := machineextensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName", "extensionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MachineExtensionsClient.List`

```go
ctx := context.TODO()
id := machineextensions.NewMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName")

// alternatively `client.List(ctx, id, machineextensions.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, machineextensions.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MachineExtensionsClient.Update`

```go
ctx := context.TODO()
id := machineextensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName", "extensionName")

payload := machineextensions.MachineExtensionUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
