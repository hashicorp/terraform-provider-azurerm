
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesetextensions` Documentation

The `virtualmachinescalesetextensions` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesetextensions"
```


### Client Initialization

```go
client := virtualmachinescalesetextensions.NewVirtualMachineScaleSetExtensionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualMachineScaleSetExtensionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualmachinescalesetextensions.NewVirtualMachineScaleSetExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "extensionName")

payload := virtualmachinescalesetextensions.VirtualMachineScaleSetExtension{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetExtensionsClient.Delete`

```go
ctx := context.TODO()
id := virtualmachinescalesetextensions.NewVirtualMachineScaleSetExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "extensionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetExtensionsClient.Get`

```go
ctx := context.TODO()
id := virtualmachinescalesetextensions.NewVirtualMachineScaleSetExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "extensionName")

read, err := client.Get(ctx, id, virtualmachinescalesetextensions.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineScaleSetExtensionsClient.List`

```go
ctx := context.TODO()
id := virtualmachinescalesetextensions.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineScaleSetExtensionsClient.Update`

```go
ctx := context.TODO()
id := virtualmachinescalesetextensions.NewVirtualMachineScaleSetExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "extensionName")

payload := virtualmachinescalesetextensions.VirtualMachineScaleSetExtensionUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
