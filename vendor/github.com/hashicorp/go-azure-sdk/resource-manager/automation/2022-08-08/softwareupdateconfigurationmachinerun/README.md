
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/softwareupdateconfigurationmachinerun` Documentation

The `softwareupdateconfigurationmachinerun` SDK allows for interaction with the Azure Resource Manager Service `automation` (API Version `2022-08-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/softwareupdateconfigurationmachinerun"
```


### Client Initialization

```go
client := softwareupdateconfigurationmachinerun.NewSoftwareUpdateConfigurationMachineRunClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SoftwareUpdateConfigurationMachineRunClient.SoftwareUpdateConfigurationMachineRunsGetById`

```go
ctx := context.TODO()
id := softwareupdateconfigurationmachinerun.NewSoftwareUpdateConfigurationMachineRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "softwareUpdateConfigurationMachineRunIdValue")

read, err := client.SoftwareUpdateConfigurationMachineRunsGetById(ctx, id, softwareupdateconfigurationmachinerun.DefaultSoftwareUpdateConfigurationMachineRunsGetByIdOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SoftwareUpdateConfigurationMachineRunClient.SoftwareUpdateConfigurationMachineRunsList`

```go
ctx := context.TODO()
id := softwareupdateconfigurationmachinerun.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue")

read, err := client.SoftwareUpdateConfigurationMachineRunsList(ctx, id, softwareupdateconfigurationmachinerun.DefaultSoftwareUpdateConfigurationMachineRunsListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
