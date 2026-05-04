
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesetrollingupgrades` Documentation

The `virtualmachinescalesetrollingupgrades` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesetrollingupgrades"
```


### Client Initialization

```go
client := virtualmachinescalesetrollingupgrades.NewVirtualMachineScaleSetRollingUpgradesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualMachineScaleSetRollingUpgradesClient.Cancel`

```go
ctx := context.TODO()
id := virtualmachinescalesetrollingupgrades.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

if err := client.CancelThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetRollingUpgradesClient.GetLatest`

```go
ctx := context.TODO()
id := virtualmachinescalesetrollingupgrades.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

read, err := client.GetLatest(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineScaleSetRollingUpgradesClient.StartExtensionUpgrade`

```go
ctx := context.TODO()
id := virtualmachinescalesetrollingupgrades.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

if err := client.StartExtensionUpgradeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetRollingUpgradesClient.StartOSUpgrade`

```go
ctx := context.TODO()
id := virtualmachinescalesetrollingupgrades.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

if err := client.StartOSUpgradeThenPoll(ctx, id); err != nil {
	// handle the error
}
```
