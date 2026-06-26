
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/deletedkeys` Documentation

The `deletedkeys` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `7.4`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/deletedkeys"
```


### Client Initialization

```go
client := deletedkeys.NewDeletedKeysClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeletedKeysClient.GetDeletedKey`

```go
ctx := context.TODO()
id := deletedkeys.NewDeletedkeyID("deletedkeyName")

read, err := client.GetDeletedKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeletedKeysClient.GetDeletedKeys`

```go
ctx := context.TODO()


// alternatively `client.GetDeletedKeys(ctx, deletedkeys.DefaultGetDeletedKeysOperationOptions())` can be used to do batched pagination
items, err := client.GetDeletedKeysComplete(ctx, deletedkeys.DefaultGetDeletedKeysOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeletedKeysClient.PurgeDeletedKey`

```go
ctx := context.TODO()
id := deletedkeys.NewDeletedkeyID("deletedkeyName")

read, err := client.PurgeDeletedKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeletedKeysClient.RecoverDeletedKey`

```go
ctx := context.TODO()
id := deletedkeys.NewDeletedkeyID("deletedkeyName")

read, err := client.RecoverDeletedKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
