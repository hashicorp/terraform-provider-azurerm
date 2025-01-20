
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesetvms` Documentation

The `virtualmachinescalesetvms` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesetvms"
```


### Client Initialization

```go
client := virtualmachinescalesetvms.NewVirtualMachineScaleSetVMsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualMachineScaleSetVMsClient.ApproveRollingUpgrade`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

if err := client.ApproveRollingUpgradeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.AttachDetachDataDisks`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

payload := virtualmachinescalesetvms.AttachDetachDataDisksRequest{
	// ...
}


if err := client.AttachDetachDataDisksThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.Deallocate`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

if err := client.DeallocateThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.Delete`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

if err := client.DeleteThenPoll(ctx, id, virtualmachinescalesetvms.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.Get`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

read, err := client.Get(ctx, id, virtualmachinescalesetvms.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.GetInstanceView`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

read, err := client.GetInstanceView(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.List`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

// alternatively `client.List(ctx, id, virtualmachinescalesetvms.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, virtualmachinescalesetvms.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.PerformMaintenance`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

if err := client.PerformMaintenanceThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.PowerOff`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

if err := client.PowerOffThenPoll(ctx, id, virtualmachinescalesetvms.DefaultPowerOffOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.Redeploy`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

if err := client.RedeployThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.Reimage`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

payload := virtualmachinescalesetvms.VirtualMachineScaleSetVMReimageParameters{
	// ...
}


if err := client.ReimageThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.ReimageAll`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

if err := client.ReimageAllThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.Restart`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

if err := client.RestartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.RetrieveBootDiagnosticsData`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

read, err := client.RetrieveBootDiagnosticsData(ctx, id, virtualmachinescalesetvms.DefaultRetrieveBootDiagnosticsDataOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.RunCommand`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

payload := virtualmachinescalesetvms.RunCommandInput{
	// ...
}


if err := client.RunCommandThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.SimulateEviction`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

read, err := client.SimulateEviction(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.Start`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetVMsClient.Update`

```go
ctx := context.TODO()
id := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "instanceId")

payload := virtualmachinescalesetvms.VirtualMachineScaleSetVM{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload, virtualmachinescalesetvms.DefaultUpdateOperationOptions()); err != nil {
	// handle the error
}
```
