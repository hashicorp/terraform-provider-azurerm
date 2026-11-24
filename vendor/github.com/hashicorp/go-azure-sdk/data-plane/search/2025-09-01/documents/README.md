
## `github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/documents` Documentation

The `documents` SDK allows for interaction with <unknown source data type> `search` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/documents"
```


### Client Initialization

```go
client := documents.NewDocumentsClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `DocumentsClient.AutocompleteGet`

```go
ctx := context.TODO()


read, err := client.AutocompleteGet(ctx, documents.DefaultAutocompleteGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DocumentsClient.AutocompletePost`

```go
ctx := context.TODO()

payload := documents.AutocompleteRequest{
	// ...
}


read, err := client.AutocompletePost(ctx, payload, documents.DefaultAutocompletePostOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DocumentsClient.Count`

```go
ctx := context.TODO()


read, err := client.Count(ctx, documents.DefaultCountOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DocumentsClient.Get`

```go
ctx := context.TODO()
id := documents.NewDocID("docName")

read, err := client.Get(ctx, id, documents.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DocumentsClient.Index`

```go
ctx := context.TODO()

payload := documents.IndexBatch{
	// ...
}


read, err := client.Index(ctx, payload, documents.DefaultIndexOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DocumentsClient.SearchGet`

```go
ctx := context.TODO()


read, err := client.SearchGet(ctx, documents.DefaultSearchGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DocumentsClient.SearchPost`

```go
ctx := context.TODO()

payload := documents.SearchRequest{
	// ...
}


read, err := client.SearchPost(ctx, payload, documents.DefaultSearchPostOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DocumentsClient.SuggestGet`

```go
ctx := context.TODO()


read, err := client.SuggestGet(ctx, documents.DefaultSuggestGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DocumentsClient.SuggestPost`

```go
ctx := context.TODO()

payload := documents.SuggestRequest{
	// ...
}


read, err := client.SuggestPost(ctx, payload, documents.DefaultSuggestPostOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
