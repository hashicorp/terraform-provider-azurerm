
## `github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachinetemplates` Documentation

The `virtualmachinetemplates` SDK allows for interaction with the Azure Resource Manager Service `systemcentervirtualmachinemanager` (API Version `2023-10-07`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachinetemplates"
```


### Client Initialization

```go
client := virtualmachinetemplates.NewVirtualMachineTemplatesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualMachineTemplatesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualmachinetemplates.NewVirtualMachineTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineTemplateValue")

payload := virtualmachinetemplates.VirtualMachineTemplate{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineTemplatesClient.Delete`

```go
ctx := context.TODO()
id := virtualmachinetemplates.NewVirtualMachineTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineTemplateValue")

if err := client.DeleteThenPoll(ctx, id, virtualmachinetemplates.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualMachineTemplatesClient.Get`

```go
ctx := context.TODO()
id := virtualmachinetemplates.NewVirtualMachineTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineTemplateValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineTemplatesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := virtualmachinetemplates.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineTemplatesClient.ListBySubscription`

```go
ctx := context.TODO()
id := virtualmachinetemplates.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineTemplatesClient.Update`

```go
ctx := context.TODO()
id := virtualmachinetemplates.NewVirtualMachineTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineTemplateValue")

payload := virtualmachinetemplates.ResourcePatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
