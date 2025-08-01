
## `github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachinegroups` Documentation

The `sqlvirtualmachinegroups` SDK allows for interaction with Azure Resource Manager `sqlvirtualmachine` (API Version `2023-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachinegroups"
```


### Client Initialization

```go
client := sqlvirtualmachinegroups.NewSqlVirtualMachineGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SqlVirtualMachineGroupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := sqlvirtualmachinegroups.NewSqlVirtualMachineGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineGroupName")

payload := sqlvirtualmachinegroups.SqlVirtualMachineGroup{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SqlVirtualMachineGroupsClient.Delete`

```go
ctx := context.TODO()
id := sqlvirtualmachinegroups.NewSqlVirtualMachineGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineGroupName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SqlVirtualMachineGroupsClient.Get`

```go
ctx := context.TODO()
id := sqlvirtualmachinegroups.NewSqlVirtualMachineGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineGroupName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SqlVirtualMachineGroupsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SqlVirtualMachineGroupsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SqlVirtualMachineGroupsClient.SqlVirtualMachinesListBySqlVMGroup`

```go
ctx := context.TODO()
id := sqlvirtualmachinegroups.NewSqlVirtualMachineGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineGroupName")

// alternatively `client.SqlVirtualMachinesListBySqlVMGroup(ctx, id)` can be used to do batched pagination
items, err := client.SqlVirtualMachinesListBySqlVMGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SqlVirtualMachineGroupsClient.Update`

```go
ctx := context.TODO()
id := sqlvirtualmachinegroups.NewSqlVirtualMachineGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sqlVirtualMachineGroupName")

payload := sqlvirtualmachinegroups.SqlVirtualMachineGroupUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
