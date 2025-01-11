
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets` Documentation

The `virtualmachinescalesets` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2024-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets"
```


### Client Initialization

```go
client := virtualmachinescalesets.NewVirtualMachineScaleSetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualMachineScaleSetsClient.ApproveRollingUpgrade`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

payload := virtualmachinescalesets.VirtualMachineScaleSetVMInstanceIDs{
	// ...
}


if err := client.ApproveRollingUpgradeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetsClient.ConvertToSinglePlacementGroup`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

payload := virtualmachinescalesets.VMScaleSetConvertToSinglePlacementGroupInput{
	// ...
}


read, err := client.ConvertToSinglePlacementGroup(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineScaleSetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

payload := virtualmachinescalesets.VirtualMachineScaleSet{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, virtualmachinescalesets.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetsClient.Deallocate`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

payload := virtualmachinescalesets.VirtualMachineScaleSetVMInstanceIDs{
	// ...
}


if err := client.DeallocateThenPoll(ctx, id, payload, virtualmachinescalesets.DefaultDeallocateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetsClient.Delete`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

if err := client.DeleteThenPoll(ctx, id, virtualmachinescalesets.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetsClient.DeleteInstances`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

payload := virtualmachinescalesets.VirtualMachineScaleSetVMInstanceRequiredIDs{
	// ...
}


if err := client.DeleteInstancesThenPoll(ctx, id, payload, virtualmachinescalesets.DefaultDeleteInstancesOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetsClient.ForceRecoveryServiceFabricPlatformUpdateDomainWalk`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

read, err := client.ForceRecoveryServiceFabricPlatformUpdateDomainWalk(ctx, id, virtualmachinescalesets.DefaultForceRecoveryServiceFabricPlatformUpdateDomainWalkOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineScaleSetsClient.Get`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

read, err := client.Get(ctx, id, virtualmachinescalesets.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineScaleSetsClient.GetInstanceView`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

read, err := client.GetInstanceView(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineScaleSetsClient.GetOSUpgradeHistory`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

// alternatively `client.GetOSUpgradeHistory(ctx, id)` can be used to do batched pagination
items, err := client.GetOSUpgradeHistoryComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineScaleSetsClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineScaleSetsClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineScaleSetsClient.ListByLocation`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.ListByLocation(ctx, id)` can be used to do batched pagination
items, err := client.ListByLocationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineScaleSetsClient.ListSkus`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

// alternatively `client.ListSkus(ctx, id)` can be used to do batched pagination
items, err := client.ListSkusComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineScaleSetsClient.PerformMaintenance`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

payload := virtualmachinescalesets.VirtualMachineScaleSetVMInstanceIDs{
	// ...
}


if err := client.PerformMaintenanceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetsClient.PowerOff`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

payload := virtualmachinescalesets.VirtualMachineScaleSetVMInstanceIDs{
	// ...
}


if err := client.PowerOffThenPoll(ctx, id, payload, virtualmachinescalesets.DefaultPowerOffOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetsClient.Reapply`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

if err := client.ReapplyThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetsClient.Redeploy`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

payload := virtualmachinescalesets.VirtualMachineScaleSetVMInstanceIDs{
	// ...
}


if err := client.RedeployThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetsClient.Reimage`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

payload := virtualmachinescalesets.VirtualMachineScaleSetReimageParameters{
	// ...
}


if err := client.ReimageThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetsClient.ReimageAll`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

payload := virtualmachinescalesets.VirtualMachineScaleSetVMInstanceIDs{
	// ...
}


if err := client.ReimageAllThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetsClient.Restart`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

payload := virtualmachinescalesets.VirtualMachineScaleSetVMInstanceIDs{
	// ...
}


if err := client.RestartThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetsClient.SetOrchestrationServiceState`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

payload := virtualmachinescalesets.OrchestrationServiceStateInput{
	// ...
}


if err := client.SetOrchestrationServiceStateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetsClient.Start`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

payload := virtualmachinescalesets.VirtualMachineScaleSetVMInstanceIDs{
	// ...
}


if err := client.StartThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetsClient.Update`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

payload := virtualmachinescalesets.VirtualMachineScaleSetUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload, virtualmachinescalesets.DefaultUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineScaleSetsClient.UpdateInstances`

```go
ctx := context.TODO()
id := virtualmachinescalesets.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

payload := virtualmachinescalesets.VirtualMachineScaleSetVMInstanceRequiredIDs{
	// ...
}


if err := client.UpdateInstancesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
