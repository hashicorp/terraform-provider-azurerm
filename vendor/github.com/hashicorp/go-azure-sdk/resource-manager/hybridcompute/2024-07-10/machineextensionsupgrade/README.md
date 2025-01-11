
## `github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/machineextensionsupgrade` Documentation

The `machineextensionsupgrade` SDK allows for interaction with Azure Resource Manager `hybridcompute` (API Version `2024-07-10`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/machineextensionsupgrade"
```


### Client Initialization

```go
client := machineextensionsupgrade.NewMachineExtensionsUpgradeClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MachineExtensionsUpgradeClient.UpgradeExtensions`

```go
ctx := context.TODO()
id := machineextensionsupgrade.NewMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName")

payload := machineextensionsupgrade.MachineExtensionUpgrade{
	// ...
}


if err := client.UpgradeExtensionsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
