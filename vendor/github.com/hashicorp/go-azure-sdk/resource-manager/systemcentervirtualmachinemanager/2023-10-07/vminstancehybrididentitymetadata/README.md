
## `github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vminstancehybrididentitymetadata` Documentation

The `vminstancehybrididentitymetadata` SDK allows for interaction with the Azure Resource Manager Service `systemcentervirtualmachinemanager` (API Version `2023-10-07`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vminstancehybrididentitymetadata"
```


### Client Initialization

```go
client := vminstancehybrididentitymetadata.NewVMInstanceHybridIdentityMetadataClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VMInstanceHybridIdentityMetadataClient.VirtualMachineInstanceHybridIdentityMetadataGet`

```go
ctx := context.TODO()
id := vminstancehybrididentitymetadata.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.VirtualMachineInstanceHybridIdentityMetadataGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VMInstanceHybridIdentityMetadataClient.VirtualMachineInstanceHybridIdentityMetadataList`

```go
ctx := context.TODO()
id := vminstancehybrididentitymetadata.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.VirtualMachineInstanceHybridIdentityMetadataList(ctx, id)` can be used to do batched pagination
items, err := client.VirtualMachineInstanceHybridIdentityMetadataListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
