
## `github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/skillsets` Documentation

The `skillsets` SDK allows for interaction with <unknown source data type> `search` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/skillsets"
```


### Client Initialization

```go
client := skillsets.NewSkillsetsClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `SkillsetsClient.Create`

```go
ctx := context.TODO()

payload := skillsets.SearchIndexerSkillset{
	// ...
}


read, err := client.Create(ctx, payload, skillsets.DefaultCreateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SkillsetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := skillsets.NewSkillsetID("skillsetName")

payload := skillsets.SearchIndexerSkillset{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, skillsets.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SkillsetsClient.Delete`

```go
ctx := context.TODO()
id := skillsets.NewSkillsetID("skillsetName")

read, err := client.Delete(ctx, id, skillsets.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SkillsetsClient.Get`

```go
ctx := context.TODO()
id := skillsets.NewSkillsetID("skillsetName")

read, err := client.Get(ctx, id, skillsets.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SkillsetsClient.List`

```go
ctx := context.TODO()


read, err := client.List(ctx, skillsets.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
