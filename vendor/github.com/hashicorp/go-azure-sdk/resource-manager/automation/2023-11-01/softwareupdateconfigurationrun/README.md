
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/softwareupdateconfigurationrun` Documentation

The `softwareupdateconfigurationrun` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/softwareupdateconfigurationrun"
```


### Client Initialization

```go
client := softwareupdateconfigurationrun.NewSoftwareUpdateConfigurationRunClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SoftwareUpdateConfigurationRunClient.GetById`

```go
ctx := context.TODO()
id := softwareupdateconfigurationrun.NewSoftwareUpdateConfigurationRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "softwareUpdateConfigurationRunId")

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
id := softwareupdateconfigurationrun.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

// alternatively `client.List(ctx, id, softwareupdateconfigurationrun.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, softwareupdateconfigurationrun.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
