
## `github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/machinenetworkprofile` Documentation

The `machinenetworkprofile` SDK allows for interaction with Azure Resource Manager `hybridcompute` (API Version `2024-07-10`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/machinenetworkprofile"
```


### Client Initialization

```go
client := machinenetworkprofile.NewMachineNetworkProfileClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MachineNetworkProfileClient.NetworkProfileGet`

```go
ctx := context.TODO()
id := machinenetworkprofile.NewMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName")

read, err := client.NetworkProfileGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
