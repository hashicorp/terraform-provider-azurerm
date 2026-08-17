
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/fullrestore` Documentation

The `fullrestore` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `7.4`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/fullrestore"
```


### Client Initialization

```go
client := fullrestore.NewFullRestoreClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `FullRestoreClient.Operation`

```go
ctx := context.TODO()

payload := fullrestore.RestoreOperationParameters{
	// ...
}


if err := client.OperationThenPoll(ctx, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FullRestoreClient.RestoreStatus`

```go
ctx := context.TODO()
id := fullrestore.NewRestoreID("jobId")

read, err := client.RestoreStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
