
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/maintenances` Documentation

The `maintenances` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2023-12-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/maintenances"
```


### Client Initialization

```go
client := maintenances.NewMaintenancesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MaintenancesClient.List`

```go
ctx := context.TODO()
id := maintenances.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MaintenancesClient.Read`

```go
ctx := context.TODO()
id := maintenances.NewMaintenanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "maintenanceName")

read, err := client.Read(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MaintenancesClient.Update`

```go
ctx := context.TODO()
id := maintenances.NewMaintenanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "maintenanceName")

payload := maintenances.MaintenanceUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
