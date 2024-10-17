
## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/datastore` Documentation

The `datastore` SDK allows for interaction with Azure Resource Manager `machinelearningservices` (API Version `2024-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/datastore"
```


### Client Initialization

```go
client := datastore.NewDatastoreClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DatastoreClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := datastore.NewDataStoreID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "dataStoreName")

payload := datastore.DatastoreResource{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, datastore.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatastoreClient.Delete`

```go
ctx := context.TODO()
id := datastore.NewDataStoreID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "dataStoreName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatastoreClient.Get`

```go
ctx := context.TODO()
id := datastore.NewDataStoreID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "dataStoreName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatastoreClient.List`

```go
ctx := context.TODO()
id := datastore.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

// alternatively `client.List(ctx, id, datastore.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, datastore.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DatastoreClient.ListSecrets`

```go
ctx := context.TODO()
id := datastore.NewDataStoreID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "dataStoreName")

read, err := client.ListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
