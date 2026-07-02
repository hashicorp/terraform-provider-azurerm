
## `github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/indexes` Documentation

The `indexes` SDK allows for interaction with <unknown source data type> `search` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/indexes"
```


### Client Initialization

```go
client := indexes.NewIndexesClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `IndexesClient.Analyze`

```go
ctx := context.TODO()
id := indexes.NewIndexID("indexName")

payload := indexes.AnalyzeRequest{
	// ...
}


read, err := client.Analyze(ctx, id, payload, indexes.DefaultAnalyzeOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IndexesClient.Create`

```go
ctx := context.TODO()

payload := indexes.SearchIndex{
	// ...
}


read, err := client.Create(ctx, payload, indexes.DefaultCreateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IndexesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := indexes.NewIndexID("indexName")

payload := indexes.SearchIndex{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, indexes.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IndexesClient.Delete`

```go
ctx := context.TODO()
id := indexes.NewIndexID("indexName")

read, err := client.Delete(ctx, id, indexes.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IndexesClient.Get`

```go
ctx := context.TODO()
id := indexes.NewIndexID("indexName")

read, err := client.Get(ctx, id, indexes.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IndexesClient.GetStatistics`

```go
ctx := context.TODO()
id := indexes.NewIndexID("indexName")

read, err := client.GetStatistics(ctx, id, indexes.DefaultGetStatisticsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IndexesClient.List`

```go
ctx := context.TODO()


read, err := client.List(ctx, indexes.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
