
## `github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/databaseprincipalassignments` Documentation

The `databaseprincipalassignments` SDK allows for interaction with Azure Resource Manager `kusto` (API Version `2023-08-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/databaseprincipalassignments"
```


### Client Initialization

```go
client := databaseprincipalassignments.NewDatabasePrincipalAssignmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DatabasePrincipalAssignmentsClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := commonids.NewKustoDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "databaseName")

payload := databaseprincipalassignments.DatabasePrincipalAssignmentCheckNameRequest{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasePrincipalAssignmentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := databaseprincipalassignments.NewDatabasePrincipalAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "databaseName", "principalAssignmentName")

payload := databaseprincipalassignments.DatabasePrincipalAssignment{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasePrincipalAssignmentsClient.Delete`

```go
ctx := context.TODO()
id := databaseprincipalassignments.NewDatabasePrincipalAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "databaseName", "principalAssignmentName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasePrincipalAssignmentsClient.Get`

```go
ctx := context.TODO()
id := databaseprincipalassignments.NewDatabasePrincipalAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "databaseName", "principalAssignmentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasePrincipalAssignmentsClient.List`

```go
ctx := context.TODO()
id := commonids.NewKustoDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "databaseName")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
