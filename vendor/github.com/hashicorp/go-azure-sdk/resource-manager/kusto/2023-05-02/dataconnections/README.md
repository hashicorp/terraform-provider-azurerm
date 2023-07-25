
## `github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/dataconnections` Documentation

The `dataconnections` SDK allows for interaction with the Azure Resource Manager Service `kusto` (API Version `2023-05-02`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/dataconnections"
```


### Client Initialization

```go
client := dataconnections.NewDataConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataConnectionsClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := dataconnections.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := dataconnections.DataConnectionCheckNameRequest{
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


### Example Usage: `DataConnectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dataconnections.NewDataConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue", "dataConnectionValue")

payload := dataconnections.DataConnection{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DataConnectionsClient.DataConnectionValidation`

```go
ctx := context.TODO()
id := dataconnections.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := dataconnections.DataConnectionValidation{
	// ...
}


if err := client.DataConnectionValidationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DataConnectionsClient.Delete`

```go
ctx := context.TODO()
id := dataconnections.NewDataConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue", "dataConnectionValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DataConnectionsClient.Get`

```go
ctx := context.TODO()
id := dataconnections.NewDataConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue", "dataConnectionValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataConnectionsClient.ListByDatabase`

```go
ctx := context.TODO()
id := dataconnections.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

read, err := client.ListByDatabase(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataConnectionsClient.Update`

```go
ctx := context.TODO()
id := dataconnections.NewDataConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue", "dataConnectionValue")

payload := dataconnections.DataConnection{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
