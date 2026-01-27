
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/secrets` Documentation

The `secrets` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/secrets"
```


### Client Initialization

```go
client := secrets.NewSecretsClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `SecretsClient.BackupSecret`

```go
ctx := context.TODO()
id := secrets.NewSecretID("secretName")

read, err := client.BackupSecret(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecretsClient.DeleteSecret`

```go
ctx := context.TODO()
id := secrets.NewSecretID("secretName")

read, err := client.DeleteSecret(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecretsClient.GetDeletedSecret`

```go
ctx := context.TODO()
id := secrets.NewDeletedsecretID("deletedsecretName")

read, err := client.GetDeletedSecret(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecretsClient.GetDeletedSecrets`

```go
ctx := context.TODO()


// alternatively `client.GetDeletedSecrets(ctx, secrets.DefaultGetDeletedSecretsOperationOptions())` can be used to do batched pagination
items, err := client.GetDeletedSecretsComplete(ctx, secrets.DefaultGetDeletedSecretsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SecretsClient.GetSecret`

```go
ctx := context.TODO()
id := secrets.NewSecretversionID("secretName", "secretversion")

read, err := client.GetSecret(ctx, id, secrets.DefaultGetSecretOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecretsClient.GetSecretVersions`

```go
ctx := context.TODO()
id := secrets.NewSecretID("secretName")

// alternatively `client.GetSecretVersions(ctx, id, secrets.DefaultGetSecretVersionsOperationOptions())` can be used to do batched pagination
items, err := client.GetSecretVersionsComplete(ctx, id, secrets.DefaultGetSecretVersionsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SecretsClient.GetSecrets`

```go
ctx := context.TODO()


// alternatively `client.GetSecrets(ctx, secrets.DefaultGetSecretsOperationOptions())` can be used to do batched pagination
items, err := client.GetSecretsComplete(ctx, secrets.DefaultGetSecretsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SecretsClient.PurgeDeletedSecret`

```go
ctx := context.TODO()
id := secrets.NewDeletedsecretID("deletedsecretName")

read, err := client.PurgeDeletedSecret(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecretsClient.RecoverDeletedSecret`

```go
ctx := context.TODO()
id := secrets.NewDeletedsecretID("deletedsecretName")

read, err := client.RecoverDeletedSecret(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecretsClient.RestoreSecret`

```go
ctx := context.TODO()

payload := secrets.SecretRestoreParameters{
	// ...
}


read, err := client.RestoreSecret(ctx, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecretsClient.SetSecret`

```go
ctx := context.TODO()
id := secrets.NewSecretID("secretName")

payload := secrets.SecretSetParameters{
	// ...
}


read, err := client.SetSecret(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecretsClient.UpdateSecret`

```go
ctx := context.TODO()
id := secrets.NewSecretversionID("secretName", "secretversion")

payload := secrets.SecretUpdateParameters{
	// ...
}


read, err := client.UpdateSecret(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
