
## `github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/datasets` Documentation

The `datasets` SDK allows for interaction with Azure Resource Manager `datafactory` (API Version `2018-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/datasets"
```


### Client Initialization

```go
client := datasets.NewDatasetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DatasetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := datasets.NewDatasetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "datasetName")

payload := datasets.DatasetResource{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, datasets.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatasetsClient.Delete`

```go
ctx := context.TODO()
id := datasets.NewDatasetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "datasetName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatasetsClient.Get`

```go
ctx := context.TODO()
id := datasets.NewDatasetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "datasetName")

read, err := client.Get(ctx, id, datasets.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatasetsClient.ListByFactory`

```go
ctx := context.TODO()
id := datasets.NewFactoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName")

// alternatively `client.ListByFactory(ctx, id)` can be used to do batched pagination
items, err := client.ListByFactoryComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
