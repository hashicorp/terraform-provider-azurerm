
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/vmsspublicipaddresses` Documentation

The `vmsspublicipaddresses` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/vmsspublicipaddresses"
```


### Client Initialization

```go
client := vmsspublicipaddresses.NewVMSSPublicIPAddressesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VMSSPublicIPAddressesClient.PublicIPAddressesGetVirtualMachineScaleSetPublicIPAddress`

```go
ctx := context.TODO()
id := commonids.NewVirtualMachineScaleSetPublicIPAddressID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "virtualMachineName", "networkInterfaceName", "ipConfigurationName", "publicIPAddressName")

read, err := client.PublicIPAddressesGetVirtualMachineScaleSetPublicIPAddress(ctx, id, vmsspublicipaddresses.DefaultPublicIPAddressesGetVirtualMachineScaleSetPublicIPAddressOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VMSSPublicIPAddressesClient.PublicIPAddressesListVirtualMachineScaleSetPublicIPAddresses`

```go
ctx := context.TODO()
id := vmsspublicipaddresses.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

// alternatively `client.PublicIPAddressesListVirtualMachineScaleSetPublicIPAddresses(ctx, id)` can be used to do batched pagination
items, err := client.PublicIPAddressesListVirtualMachineScaleSetPublicIPAddressesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VMSSPublicIPAddressesClient.PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddresses`

```go
ctx := context.TODO()
id := commonids.NewVirtualMachineScaleSetIPConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "virtualMachineName", "networkInterfaceName", "ipConfigurationName")

// alternatively `client.PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddresses(ctx, id)` can be used to do batched pagination
items, err := client.PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddressesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
