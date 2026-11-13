
## `github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/datasources` Documentation

The `datasources` SDK allows for interaction with <unknown source data type> `search` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/datasources"
```


### Client Initialization

```go
client := datasources.NewDataSourcesClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataSourcesClient.Create`

```go
ctx := context.TODO()

payload := datasources.SearchIndexerDataSource{
	// ...
}


read, err := client.Create(ctx, payload, datasources.DefaultCreateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataSourcesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := datasources.NewDatasourceID("datasourceName")

payload := datasources.SearchIndexerDataSource{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, datasources.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataSourcesClient.Delete`

```go
ctx := context.TODO()
id := datasources.NewDatasourceID("datasourceName")

read, err := client.Delete(ctx, id, datasources.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataSourcesClient.Get`

```go
ctx := context.TODO()
id := datasources.NewDatasourceID("datasourceName")

read, err := client.Get(ctx, id, datasources.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataSourcesClient.List`

```go
ctx := context.TODO()


read, err := client.List(ctx, datasources.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
