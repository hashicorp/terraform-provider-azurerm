
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/softwareupdateconfiguration` Documentation

The `softwareupdateconfiguration` SDK allows for interaction with the Azure Resource Manager Service `automation` (API Version `2019-06-01`).

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


### Example Usage: `SoftwareUpdateConfigurationClient.SoftwareUpdateConfigurationsCreate`

```go
ctx := context.TODO()
id := softwareupdateconfiguration.NewSoftwareUpdateConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "softwareUpdateConfigurationValue")

payload := softwareupdateconfiguration.SoftwareUpdateConfiguration{
	// ...
}


read, err := client.SoftwareUpdateConfigurationsCreate(ctx, id, payload, softwareupdateconfiguration.DefaultSoftwareUpdateConfigurationsCreateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SoftwareUpdateConfigurationClient.SoftwareUpdateConfigurationsDelete`

```go
ctx := context.TODO()
id := softwareupdateconfiguration.NewSoftwareUpdateConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "softwareUpdateConfigurationValue")

read, err := client.SoftwareUpdateConfigurationsDelete(ctx, id, softwareupdateconfiguration.DefaultSoftwareUpdateConfigurationsDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SoftwareUpdateConfigurationClient.SoftwareUpdateConfigurationsGetByName`

```go
ctx := context.TODO()
id := softwareupdateconfiguration.NewSoftwareUpdateConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "softwareUpdateConfigurationValue")

read, err := client.SoftwareUpdateConfigurationsGetByName(ctx, id, softwareupdateconfiguration.DefaultSoftwareUpdateConfigurationsGetByNameOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SoftwareUpdateConfigurationClient.SoftwareUpdateConfigurationsList`

```go
ctx := context.TODO()
id := softwareupdateconfiguration.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue")

read, err := client.SoftwareUpdateConfigurationsList(ctx, id, softwareupdateconfiguration.DefaultSoftwareUpdateConfigurationsListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
