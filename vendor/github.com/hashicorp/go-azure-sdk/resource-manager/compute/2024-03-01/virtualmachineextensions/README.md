
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachineextensions` Documentation

The `virtualmachineextensions` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachineextensions"
```


### Client Initialization

```go
client := virtualmachineextensions.NewVirtualMachineExtensionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualMachineExtensionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualmachineextensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "extensionName")

payload := virtualmachineextensions.VirtualMachineExtension{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineExtensionsClient.Delete`

```go
ctx := context.TODO()
id := virtualmachineextensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "extensionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineExtensionsClient.Get`

```go
ctx := context.TODO()
id := virtualmachineextensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "extensionName")

read, err := client.Get(ctx, id, virtualmachineextensions.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineExtensionsClient.List`

```go
ctx := context.TODO()
id := virtualmachineextensions.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

read, err := client.List(ctx, id, virtualmachineextensions.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineExtensionsClient.Update`

```go
ctx := context.TODO()
id := virtualmachineextensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "extensionName")

payload := virtualmachineextensions.VirtualMachineExtensionUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
