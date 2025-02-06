
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines` Documentation

The `virtualmachines` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
```


### Client Initialization

```go
client := virtualmachines.NewVirtualMachinesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualMachinesClient.AssessPatches`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

if err := client.AssessPatchesThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.AttachDetachDataDisks`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

payload := virtualmachines.AttachDetachDataDisksRequest{
	// ...
}


if err := client.AttachDetachDataDisksThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Capture`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

payload := virtualmachines.VirtualMachineCaptureParameters{
	// ...
}


if err := client.CaptureThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.ConvertToManagedDisks`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

if err := client.ConvertToManagedDisksThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

payload := virtualmachines.VirtualMachine{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Deallocate`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

if err := client.DeallocateThenPoll(ctx, id, virtualmachines.DefaultDeallocateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Delete`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

if err := client.DeleteThenPoll(ctx, id, virtualmachines.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Generalize`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

read, err := client.Generalize(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachinesClient.Get`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

read, err := client.Get(ctx, id, virtualmachines.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachinesClient.InstallPatches`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

payload := virtualmachines.VirtualMachineInstallPatchesParameters{
	// ...
}


if err := client.InstallPatchesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.InstanceView`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

read, err := client.InstanceView(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachinesClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id, virtualmachines.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, virtualmachines.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachinesClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id, virtualmachines.DefaultListAllOperationOptions())` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id, virtualmachines.DefaultListAllOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachinesClient.ListAvailableSizes`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

read, err := client.ListAvailableSizes(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachinesClient.ListByLocation`

```go
ctx := context.TODO()
id := virtualmachines.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.ListByLocation(ctx, id)` can be used to do batched pagination
items, err := client.ListByLocationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachinesClient.PerformMaintenance`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

if err := client.PerformMaintenanceThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.PowerOff`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

if err := client.PowerOffThenPoll(ctx, id, virtualmachines.DefaultPowerOffOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Reapply`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

if err := client.ReapplyThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Redeploy`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

if err := client.RedeployThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Reimage`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

payload := virtualmachines.VirtualMachineReimageParameters{
	// ...
}


if err := client.ReimageThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Restart`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

if err := client.RestartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.RetrieveBootDiagnosticsData`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

read, err := client.RetrieveBootDiagnosticsData(ctx, id, virtualmachines.DefaultRetrieveBootDiagnosticsDataOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachinesClient.RunCommand`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

payload := virtualmachines.RunCommandInput{
	// ...
}


if err := client.RunCommandThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.SimulateEviction`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

read, err := client.SimulateEviction(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachinesClient.Start`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Update`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

payload := virtualmachines.VirtualMachineUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload, virtualmachines.DefaultUpdateOperationOptions()); err != nil {
	// handle the error
}
```
