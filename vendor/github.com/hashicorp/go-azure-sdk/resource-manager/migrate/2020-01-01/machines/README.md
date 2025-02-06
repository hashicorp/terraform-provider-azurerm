
## `github.com/hashicorp/go-azure-sdk/resource-manager/migrate/2020-01-01/machines` Documentation

The `machines` SDK allows for interaction with Azure Resource Manager `migrate` (API Version `2020-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/migrate/2020-01-01/machines"
```


### Client Initialization

```go
client := machines.NewMachinesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MachinesClient.GetAllMachinesInSite`

```go
ctx := context.TODO()
id := machines.NewVMwareSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vmwareSiteName")

// alternatively `client.GetAllMachinesInSite(ctx, id, machines.DefaultGetAllMachinesInSiteOperationOptions())` can be used to do batched pagination
items, err := client.GetAllMachinesInSiteComplete(ctx, id, machines.DefaultGetAllMachinesInSiteOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MachinesClient.GetMachine`

```go
ctx := context.TODO()
id := commonids.NewVMwareSiteMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vmwareSiteName", "machineName")

read, err := client.GetMachine(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MachinesClient.StartMachine`

```go
ctx := context.TODO()
id := commonids.NewVMwareSiteMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vmwareSiteName", "machineName")

read, err := client.StartMachine(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MachinesClient.StopMachine`

```go
ctx := context.TODO()
id := commonids.NewVMwareSiteMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vmwareSiteName", "machineName")

read, err := client.StopMachine(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
