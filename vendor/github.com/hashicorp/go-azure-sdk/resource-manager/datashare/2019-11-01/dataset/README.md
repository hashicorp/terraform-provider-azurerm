
## `github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset` Documentation

The `dataset` SDK allows for interaction with the Azure Resource Manager Service `datashare` (API Version `2019-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
```


### Client Initialization

```go
client := dataset.NewDataSetClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataSetClient.Create`

```go
ctx := context.TODO()
id := dataset.NewDataSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "shareValue", "dataSetValue")

payload := dataset.DataSet{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataSetClient.Delete`

```go
ctx := context.TODO()
id := dataset.NewDataSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "shareValue", "dataSetValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DataSetClient.Get`

```go
ctx := context.TODO()
id := dataset.NewDataSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "shareValue", "dataSetValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataSetClient.ListByShare`

```go
ctx := context.TODO()
id := dataset.NewShareID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "shareValue")

// alternatively `client.ListByShare(ctx, id, dataset.DefaultListByShareOperationOptions())` can be used to do batched pagination
items, err := client.ListByShareComplete(ctx, id, dataset.DefaultListByShareOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
