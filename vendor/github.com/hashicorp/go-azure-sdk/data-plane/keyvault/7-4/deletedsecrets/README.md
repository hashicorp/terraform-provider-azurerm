
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/deletedsecrets` Documentation

The `deletedsecrets` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `7.4`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/deletedsecrets"
```


### Client Initialization

```go
client := deletedsecrets.NewDeletedSecretsClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeletedSecretsClient.GetDeletedSecret`

```go
ctx := context.TODO()
id := deletedsecrets.NewDeletedsecretID("deletedsecretName")

read, err := client.GetDeletedSecret(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeletedSecretsClient.GetDeletedSecrets`

```go
ctx := context.TODO()


// alternatively `client.GetDeletedSecrets(ctx, deletedsecrets.DefaultGetDeletedSecretsOperationOptions())` can be used to do batched pagination
items, err := client.GetDeletedSecretsComplete(ctx, deletedsecrets.DefaultGetDeletedSecretsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeletedSecretsClient.PurgeDeletedSecret`

```go
ctx := context.TODO()
id := deletedsecrets.NewDeletedsecretID("deletedsecretName")

read, err := client.PurgeDeletedSecret(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeletedSecretsClient.RecoverDeletedSecret`

```go
ctx := context.TODO()
id := deletedsecrets.NewDeletedsecretID("deletedsecretName")

read, err := client.RecoverDeletedSecret(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
