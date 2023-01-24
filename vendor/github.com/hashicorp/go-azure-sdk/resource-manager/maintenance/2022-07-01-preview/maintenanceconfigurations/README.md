
## `github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/maintenanceconfigurations` Documentation

The `maintenanceconfigurations` SDK allows for interaction with the Azure Resource Manager Service `maintenance` (API Version `2022-07-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/maintenanceconfigurations"
```


### Client Initialization

```go
client := maintenanceconfigurations.NewMaintenanceConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MaintenanceConfigurationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := maintenanceconfigurations.NewMaintenanceConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "maintenanceConfigurationValue")

payload := maintenanceconfigurations.MaintenanceConfiguration{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MaintenanceConfigurationsClient.Delete`

```go
ctx := context.TODO()
id := maintenanceconfigurations.NewMaintenanceConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "maintenanceConfigurationValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MaintenanceConfigurationsClient.ForResourceGroupList`

```go
ctx := context.TODO()
id := maintenanceconfigurations.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.ForResourceGroupList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MaintenanceConfigurationsClient.Get`

```go
ctx := context.TODO()
id := maintenanceconfigurations.NewMaintenanceConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "maintenanceConfigurationValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MaintenanceConfigurationsClient.List`

```go
ctx := context.TODO()
id := maintenanceconfigurations.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MaintenanceConfigurationsClient.Update`

```go
ctx := context.TODO()
id := maintenanceconfigurations.NewMaintenanceConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "maintenanceConfigurationValue")

payload := maintenanceconfigurations.MaintenanceConfiguration{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
