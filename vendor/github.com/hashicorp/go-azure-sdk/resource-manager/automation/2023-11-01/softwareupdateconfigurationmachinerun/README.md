
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/softwareupdateconfigurationmachinerun` Documentation

The `softwareupdateconfigurationmachinerun` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/softwareupdateconfigurationmachinerun"
```


### Client Initialization

```go
client := softwareupdateconfigurationmachinerun.NewSoftwareUpdateConfigurationMachineRunClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SoftwareUpdateConfigurationMachineRunClient.GetById`

```go
ctx := context.TODO()
id := softwareupdateconfigurationmachinerun.NewSoftwareUpdateConfigurationMachineRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "softwareUpdateConfigurationMachineRunId")

read, err := client.GetById(ctx, id, softwareupdateconfigurationmachinerun.DefaultGetByIdOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SoftwareUpdateConfigurationMachineRunClient.List`

```go
ctx := context.TODO()
id := softwareupdateconfigurationmachinerun.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

// alternatively `client.List(ctx, id, softwareupdateconfigurationmachinerun.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, softwareupdateconfigurationmachinerun.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
