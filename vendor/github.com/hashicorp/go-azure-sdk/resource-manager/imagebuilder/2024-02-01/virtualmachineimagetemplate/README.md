
## `github.com/hashicorp/go-azure-sdk/resource-manager/imagebuilder/2024-02-01/virtualmachineimagetemplate` Documentation

The `virtualmachineimagetemplate` SDK allows for interaction with Azure Resource Manager `imagebuilder` (API Version `2024-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/imagebuilder/2024-02-01/virtualmachineimagetemplate"
```


### Client Initialization

```go
client := virtualmachineimagetemplate.NewVirtualMachineImageTemplateClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualMachineImageTemplateClient.Cancel`

```go
ctx := context.TODO()
id := virtualmachineimagetemplate.NewImageTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "imageTemplateName")

if err := client.CancelThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineImageTemplateClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualmachineimagetemplate.NewImageTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "imageTemplateName")

payload := virtualmachineimagetemplate.ImageTemplate{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineImageTemplateClient.Delete`

```go
ctx := context.TODO()
id := virtualmachineimagetemplate.NewImageTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "imageTemplateName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineImageTemplateClient.Get`

```go
ctx := context.TODO()
id := virtualmachineimagetemplate.NewImageTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "imageTemplateName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineImageTemplateClient.GetRunOutput`

```go
ctx := context.TODO()
id := virtualmachineimagetemplate.NewRunOutputID("12345678-1234-9876-4563-123456789012", "example-resource-group", "imageTemplateName", "runOutputName")

read, err := client.GetRunOutput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineImageTemplateClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineImageTemplateClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineImageTemplateClient.ListRunOutputs`

```go
ctx := context.TODO()
id := virtualmachineimagetemplate.NewImageTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "imageTemplateName")

// alternatively `client.ListRunOutputs(ctx, id)` can be used to do batched pagination
items, err := client.ListRunOutputsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineImageTemplateClient.Run`

```go
ctx := context.TODO()
id := virtualmachineimagetemplate.NewImageTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "imageTemplateName")

if err := client.RunThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineImageTemplateClient.Update`

```go
ctx := context.TODO()
id := virtualmachineimagetemplate.NewImageTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "imageTemplateName")

payload := virtualmachineimagetemplate.ImageTemplateUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
