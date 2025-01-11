
## `github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualmachines` Documentation

The `virtualmachines` SDK allows for interaction with Azure Resource Manager `devtestlab` (API Version `2018-09-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualmachines"
```


### Client Initialization

```go
client := virtualmachines.NewVirtualMachinesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualMachinesClient.AddDataDisk`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

payload := virtualmachines.DataDiskProperties{
	// ...
}


if err := client.AddDataDiskThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.ApplyArtifacts`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

payload := virtualmachines.ApplyArtifactsRequest{
	// ...
}


if err := client.ApplyArtifactsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Claim`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

if err := client.ClaimThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

payload := virtualmachines.LabVirtualMachine{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Delete`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.DetachDataDisk`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

payload := virtualmachines.DetachDataDiskProperties{
	// ...
}


if err := client.DetachDataDiskThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Get`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

read, err := client.Get(ctx, id, virtualmachines.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachinesClient.GetRdpFileContents`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

read, err := client.GetRdpFileContents(ctx, id)
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
id := virtualmachines.NewLabID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName")

// alternatively `client.List(ctx, id, virtualmachines.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, virtualmachines.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachinesClient.ListApplicableSchedules`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

read, err := client.ListApplicableSchedules(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachinesClient.Redeploy`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

if err := client.RedeployThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Resize`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

payload := virtualmachines.ResizeLabVirtualMachineProperties{
	// ...
}


if err := client.ResizeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Restart`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

if err := client.RestartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Start`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Stop`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

if err := client.StopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.TransferDisks`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

if err := client.TransferDisksThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.UnClaim`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

if err := client.UnClaimThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachinesClient.Update`

```go
ctx := context.TODO()
id := virtualmachines.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labName", "virtualMachineName")

payload := virtualmachines.UpdateResource{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
