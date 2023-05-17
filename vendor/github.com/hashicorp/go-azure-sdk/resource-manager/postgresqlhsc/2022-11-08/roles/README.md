
## `github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/roles` Documentation

The `roles` SDK allows for interaction with the Azure Resource Manager Service `postgresqlhsc` (API Version `2022-11-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/roles"
```


### Client Initialization

```go
client := roles.NewRolesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RolesClient.Create`

```go
ctx := context.TODO()
id := roles.NewRoleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverGroupsv2Value", "roleValue")

payload := roles.Role{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RolesClient.Delete`

```go
ctx := context.TODO()
id := roles.NewRoleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverGroupsv2Value", "roleValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RolesClient.Get`

```go
ctx := context.TODO()
id := roles.NewRoleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverGroupsv2Value", "roleValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RolesClient.ListByCluster`

```go
ctx := context.TODO()
id := roles.NewServerGroupsv2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverGroupsv2Value")

read, err := client.ListByCluster(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
