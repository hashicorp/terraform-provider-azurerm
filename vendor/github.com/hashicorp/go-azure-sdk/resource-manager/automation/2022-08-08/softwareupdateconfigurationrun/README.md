
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/softwareupdateconfigurationrun` Documentation

The `softwareupdateconfigurationrun` SDK allows for interaction with the Azure Resource Manager Service `automation` (API Version `2022-08-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/softwareupdateconfigurationrun"
```


### Client Initialization

```go
client := softwareupdateconfigurationrun.NewSoftwareUpdateConfigurationRunClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SoftwareUpdateConfigurationRunClient.GetById`

```go
ctx := context.TODO()
id := softwareupdateconfigurationrun.NewSoftwareUpdateConfigurationRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "softwareUpdateConfigurationRunIdValue")

read, err := client.GetById(ctx, id, softwareupdateconfigurationrun.DefaultGetByIdOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SoftwareUpdateConfigurationRunClient.List`

```go
ctx := context.TODO()
id := softwareupdateconfigurationrun.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue")

read, err := client.List(ctx, id, softwareupdateconfigurationrun.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
