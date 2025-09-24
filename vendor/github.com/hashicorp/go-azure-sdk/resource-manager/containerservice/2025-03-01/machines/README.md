
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/machines` Documentation

The `machines` SDK allows for interaction with Azure Resource Manager `containerservice` (API Version `2025-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/machines"
```


### Client Initialization

```go
client := machines.NewMachinesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MachinesClient.Get`

```go
ctx := context.TODO()
id := machines.NewMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName", "agentPoolName", "machineName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MachinesClient.List`

```go
ctx := context.TODO()
id := machines.NewAgentPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName", "agentPoolName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
