
## `github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/rbacs` Documentation

The `rbacs` SDK allows for interaction with Azure Resource Manager `cosmosdb` (API Version `2024-08-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/rbacs"
```


### Client Initialization

```go
client := rbacs.NewRbacsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RbacsClient.SqlResourcesCreateUpdateSqlRoleAssignment`

```go
ctx := context.TODO()
id := rbacs.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := rbacs.SqlRoleAssignmentCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateSqlRoleAssignmentThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RbacsClient.SqlResourcesCreateUpdateSqlRoleDefinition`

```go
ctx := context.TODO()
id := rbacs.NewSqlRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

payload := rbacs.SqlRoleDefinitionCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateSqlRoleDefinitionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RbacsClient.SqlResourcesDeleteSqlRoleAssignment`

```go
ctx := context.TODO()
id := rbacs.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

if err := client.SqlResourcesDeleteSqlRoleAssignmentThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RbacsClient.SqlResourcesDeleteSqlRoleDefinition`

```go
ctx := context.TODO()
id := rbacs.NewSqlRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

if err := client.SqlResourcesDeleteSqlRoleDefinitionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RbacsClient.SqlResourcesGetSqlRoleAssignment`

```go
ctx := context.TODO()
id := rbacs.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.SqlResourcesGetSqlRoleAssignment(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RbacsClient.SqlResourcesGetSqlRoleDefinition`

```go
ctx := context.TODO()
id := rbacs.NewSqlRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

read, err := client.SqlResourcesGetSqlRoleDefinition(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RbacsClient.SqlResourcesListSqlRoleAssignments`

```go
ctx := context.TODO()
id := rbacs.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

read, err := client.SqlResourcesListSqlRoleAssignments(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RbacsClient.SqlResourcesListSqlRoleDefinitions`

```go
ctx := context.TODO()
id := rbacs.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

read, err := client.SqlResourcesListSqlRoleDefinitions(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
