
## `github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/synonymmaps` Documentation

The `synonymmaps` SDK allows for interaction with <unknown source data type> `search` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/synonymmaps"
```


### Client Initialization

```go
client := synonymmaps.NewSynonymMapsClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `SynonymMapsClient.Create`

```go
ctx := context.TODO()

payload := synonymmaps.SynonymMap{
	// ...
}


read, err := client.Create(ctx, payload, synonymmaps.DefaultCreateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SynonymMapsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := synonymmaps.NewSynonymmapID("synonymmapName")

payload := synonymmaps.SynonymMap{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, synonymmaps.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SynonymMapsClient.Delete`

```go
ctx := context.TODO()
id := synonymmaps.NewSynonymmapID("synonymmapName")

read, err := client.Delete(ctx, id, synonymmaps.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SynonymMapsClient.Get`

```go
ctx := context.TODO()
id := synonymmaps.NewSynonymmapID("synonymmapName")

read, err := client.Get(ctx, id, synonymmaps.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SynonymMapsClient.List`

```go
ctx := context.TODO()


read, err := client.List(ctx, synonymmaps.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
