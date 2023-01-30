
## `github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/publicmaintenanceconfigurations` Documentation

The `publicmaintenanceconfigurations` SDK allows for interaction with the Azure Resource Manager Service `maintenance` (API Version `2022-07-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/publicmaintenanceconfigurations"
```


### Client Initialization

```go
client := publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PublicMaintenanceConfigurationsClient.Get`

```go
ctx := context.TODO()
id := publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID("12345678-1234-9876-4563-123456789012", "publicMaintenanceConfigurationValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PublicMaintenanceConfigurationsClient.List`

```go
ctx := context.TODO()
id := publicmaintenanceconfigurations.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
