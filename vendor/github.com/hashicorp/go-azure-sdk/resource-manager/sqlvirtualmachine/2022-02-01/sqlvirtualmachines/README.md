
## `github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/sqlvirtualmachines` Documentation

The `sqlvirtualmachines` SDK allows for interaction with the Azure Resource Manager Service `sqlvirtualmachine` (API Version `2022-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/sqlvirtualmachines"
```


### Client Initialization

```go
client := sqlvirtualmachines.NewSqlVirtualMachinesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
if err != nil {
	// handle the error
}
```


### Example Usage: `SqlVirtualMachinesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := sqlvirtualmachines.NewSqlVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineValue")

payload := sqlvirtualmachines.SqlVirtualMachine{
	// ...
}

future, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `SqlVirtualMachinesClient.Delete`

```go
ctx := context.TODO()
id := sqlvirtualmachines.NewSqlVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineValue")
future, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `SqlVirtualMachinesClient.Get`

```go
ctx := context.TODO()
id := sqlvirtualmachines.NewSqlVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineValue")
read, err := client.Get(ctx, id, sqlvirtualmachines.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SqlVirtualMachinesClient.List`

```go
ctx := context.TODO()
id := sqlvirtualmachines.NewSubscriptionID()
// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SqlVirtualMachinesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := sqlvirtualmachines.NewResourceGroupID()
// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SqlVirtualMachinesClient.ListBySqlVmGroup`

```go
ctx := context.TODO()
id := sqlvirtualmachines.NewSqlVirtualMachineGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineGroupValue")
// alternatively `client.ListBySqlVmGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListBySqlVmGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SqlVirtualMachinesClient.Redeploy`

```go
ctx := context.TODO()
id := sqlvirtualmachines.NewSqlVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineValue")
future, err := client.Redeploy(ctx, id)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `SqlVirtualMachinesClient.StartAssessment`

```go
ctx := context.TODO()
id := sqlvirtualmachines.NewSqlVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineValue")
future, err := client.StartAssessment(ctx, id)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `SqlVirtualMachinesClient.Update`

```go
ctx := context.TODO()
id := sqlvirtualmachines.NewSqlVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineValue")

payload := sqlvirtualmachines.SqlVirtualMachineUpdate{
	// ...
}

future, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```
