
## `github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-05-20-preview/machineextensionsupgrade` Documentation

The `machineextensionsupgrade` SDK allows for interaction with the Azure Resource Manager Service `hybridcompute` (API Version `2024-05-20-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-05-20-preview/machineextensionsupgrade"
```


### Client Initialization

```go
client := machineextensionsupgrade.NewMachineExtensionsUpgradeClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MachineExtensionsUpgradeClient.UpgradeExtensions`

```go
ctx := context.TODO()
id := machineextensionsupgrade.NewMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineValue")

payload := machineextensionsupgrade.MachineExtensionUpgrade{
	// ...
}


if err := client.UpgradeExtensionsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
