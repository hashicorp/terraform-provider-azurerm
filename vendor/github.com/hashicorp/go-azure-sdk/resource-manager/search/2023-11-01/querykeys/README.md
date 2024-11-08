
## `github.com/hashicorp/go-azure-sdk/resource-manager/search/2023-11-01/querykeys` Documentation

The `querykeys` SDK allows for interaction with Azure Resource Manager `search` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/search/2023-11-01/querykeys"
```


### Client Initialization

```go
client := querykeys.NewQueryKeysClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `QueryKeysClient.Create`

```go
ctx := context.TODO()
id := querykeys.NewCreateQueryKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceName", "createQueryKeyName")

read, err := client.Create(ctx, id, querykeys.DefaultCreateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueryKeysClient.Delete`

```go
ctx := context.TODO()
id := querykeys.NewDeleteQueryKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceName", "deleteQueryKeyName")

read, err := client.Delete(ctx, id, querykeys.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueryKeysClient.ListBySearchService`

```go
ctx := context.TODO()
id := querykeys.NewSearchServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceName")

// alternatively `client.ListBySearchService(ctx, id, querykeys.DefaultListBySearchServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListBySearchServiceComplete(ctx, id, querykeys.DefaultListBySearchServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
