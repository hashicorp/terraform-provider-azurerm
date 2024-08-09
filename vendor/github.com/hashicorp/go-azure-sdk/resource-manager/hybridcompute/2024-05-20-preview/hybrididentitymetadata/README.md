
## `github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-05-20-preview/hybrididentitymetadata` Documentation

The `hybrididentitymetadata` SDK allows for interaction with the Azure Resource Manager Service `hybridcompute` (API Version `2024-05-20-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-05-20-preview/hybrididentitymetadata"
```


### Client Initialization

```go
client := hybrididentitymetadata.NewHybridIdentityMetadataClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `HybridIdentityMetadataClient.Get`

```go
ctx := context.TODO()
id := hybrididentitymetadata.NewHybridIdentityMetadataID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineValue", "hybridIdentityMetadataValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HybridIdentityMetadataClient.ListByMachines`

```go
ctx := context.TODO()
id := hybrididentitymetadata.NewMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineValue")

// alternatively `client.ListByMachines(ctx, id)` can be used to do batched pagination
items, err := client.ListByMachinesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
