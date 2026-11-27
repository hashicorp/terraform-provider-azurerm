
## `github.com/hashicorp/go-azure-sdk/resource-manager/resourcegraph/2024-04-01/graphqueries` Documentation

The `graphqueries` SDK allows for interaction with Azure Resource Manager `resourcegraph` (API Version `2024-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/resourcegraph/2024-04-01/graphqueries"
```


### Client Initialization

```go
client := graphqueries.NewGraphqueriesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GraphqueriesClient.GraphQueryUpdate`

```go
ctx := context.TODO()
id := graphqueries.NewQueryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "queryName")

payload := graphqueries.GraphQueryUpdateParameters{
	// ...
}


read, err := client.GraphQueryUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
