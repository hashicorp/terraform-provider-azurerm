
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/operationprogress` Documentation

The `operationprogress` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2023-12-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/operationprogress"
```


### Client Initialization

```go
client := operationprogress.NewOperationProgressClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OperationProgressClient.Get`

```go
ctx := context.TODO()
id := operationprogress.NewOperationProgressID("12345678-1234-9876-4563-123456789012", "locationName", "operationId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
