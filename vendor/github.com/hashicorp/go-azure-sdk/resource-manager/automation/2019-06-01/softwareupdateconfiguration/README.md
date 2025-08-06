
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/softwareupdateconfiguration` Documentation

The `softwareupdateconfiguration` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2019-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/softwareupdateconfiguration"
```


### Client Initialization

```go
client := softwareupdateconfiguration.NewSoftwareUpdateConfigurationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SoftwareUpdateConfigurationClient.Create`

```go
ctx := context.TODO()
id := softwareupdateconfiguration.NewSoftwareUpdateConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "softwareUpdateConfigurationName")

payload := softwareupdateconfiguration.SoftwareUpdateConfiguration{
	// ...
}


read, err := client.Create(ctx, id, payload, softwareupdateconfiguration.DefaultCreateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SoftwareUpdateConfigurationClient.Delete`

```go
ctx := context.TODO()
id := softwareupdateconfiguration.NewSoftwareUpdateConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "softwareUpdateConfigurationName")

read, err := client.Delete(ctx, id, softwareupdateconfiguration.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SoftwareUpdateConfigurationClient.GetByName`

```go
ctx := context.TODO()
id := softwareupdateconfiguration.NewSoftwareUpdateConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "softwareUpdateConfigurationName")

read, err := client.GetByName(ctx, id, softwareupdateconfiguration.DefaultGetByNameOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SoftwareUpdateConfigurationClient.List`

```go
ctx := context.TODO()
id := softwareupdateconfiguration.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

read, err := client.List(ctx, id, softwareupdateconfiguration.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
