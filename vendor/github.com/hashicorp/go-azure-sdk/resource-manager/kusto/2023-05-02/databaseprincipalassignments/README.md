
## `github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/databaseprincipalassignments` Documentation

The `databaseprincipalassignments` SDK allows for interaction with the Azure Resource Manager Service `kusto` (API Version `2023-05-02`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/databaseprincipalassignments"
```


### Client Initialization

```go
client := databaseprincipalassignments.NewDatabasePrincipalAssignmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DatabasePrincipalAssignmentsClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := databaseprincipalassignments.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

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
id := databaseprincipalassignments.NewDatabasePrincipalAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue", "principalAssignmentValue")

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
id := databaseprincipalassignments.NewDatabasePrincipalAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue", "principalAssignmentValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasePrincipalAssignmentsClient.Get`

```go
ctx := context.TODO()
id := databaseprincipalassignments.NewDatabasePrincipalAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue", "principalAssignmentValue")

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
id := databaseprincipalassignments.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
