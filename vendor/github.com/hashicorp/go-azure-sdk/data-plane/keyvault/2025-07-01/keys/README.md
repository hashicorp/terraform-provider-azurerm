
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/keys` Documentation

The `keys` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/keys"
```


### Client Initialization

```go
client := keys.NewKeysClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `KeysClient.BackupKey`

```go
ctx := context.TODO()
id := keys.NewKeyID("keyName")

read, err := client.BackupKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.CreateKey`

```go
ctx := context.TODO()
id := keys.NewKeyID("keyName")

payload := keys.KeyCreateParameters{
	// ...
}


read, err := client.CreateKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.Decrypt`

```go
ctx := context.TODO()
id := keys.NewKeyversionID("https://endpoint-url.example.com", "keyName")

payload := keys.KeyOperationsParameters{
	// ...
}


read, err := client.Decrypt(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.DeleteKey`

```go
ctx := context.TODO()
id := keys.NewKeyID("keyName")

read, err := client.DeleteKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.Encrypt`

```go
ctx := context.TODO()
id := keys.NewKeyversionID("https://endpoint-url.example.com", "keyName")

payload := keys.KeyOperationsParameters{
	// ...
}


read, err := client.Encrypt(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.GetDeletedKey`

```go
ctx := context.TODO()
id := keys.NewDeletedkeyID("deletedkeyName")

read, err := client.GetDeletedKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.GetDeletedKeys`

```go
ctx := context.TODO()


// alternatively `client.GetDeletedKeys(ctx, keys.DefaultGetDeletedKeysOperationOptions())` can be used to do batched pagination
items, err := client.GetDeletedKeysComplete(ctx, keys.DefaultGetDeletedKeysOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `KeysClient.GetKey`

```go
ctx := context.TODO()
id := keys.NewKeyversionID("https://endpoint-url.example.com", "keyName")

read, err := client.GetKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.GetKeyAttestation`

```go
ctx := context.TODO()
id := keys.NewKeyversionID("https://endpoint-url.example.com", "keyName")

read, err := client.GetKeyAttestation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.GetKeyRotationPolicy`

```go
ctx := context.TODO()
id := keys.NewKeyID("keyName")

read, err := client.GetKeyRotationPolicy(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.GetKeyVersions`

```go
ctx := context.TODO()
id := keys.NewKeyID("keyName")

// alternatively `client.GetKeyVersions(ctx, id, keys.DefaultGetKeyVersionsOperationOptions())` can be used to do batched pagination
items, err := client.GetKeyVersionsComplete(ctx, id, keys.DefaultGetKeyVersionsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `KeysClient.GetKeys`

```go
ctx := context.TODO()


// alternatively `client.GetKeys(ctx, keys.DefaultGetKeysOperationOptions())` can be used to do batched pagination
items, err := client.GetKeysComplete(ctx, keys.DefaultGetKeysOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `KeysClient.GetRandomBytes`

```go
ctx := context.TODO()

payload := keys.GetRandomBytesRequest{
	// ...
}


read, err := client.GetRandomBytes(ctx, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.Ign`

```go
ctx := context.TODO()
id := keys.NewKeyversionID("https://endpoint-url.example.com", "keyName")

payload := keys.KeySignParameters{
	// ...
}


read, err := client.Ign(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.ImportKey`

```go
ctx := context.TODO()
id := keys.NewKeyID("keyName")

payload := keys.KeyImportParameters{
	// ...
}


read, err := client.ImportKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.PurgeDeletedKey`

```go
ctx := context.TODO()
id := keys.NewDeletedkeyID("deletedkeyName")

read, err := client.PurgeDeletedKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.RecoverDeletedKey`

```go
ctx := context.TODO()
id := keys.NewDeletedkeyID("deletedkeyName")

read, err := client.RecoverDeletedKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.Release`

```go
ctx := context.TODO()
id := keys.NewKeyversionID("https://endpoint-url.example.com", "keyName")

payload := keys.KeyReleaseParameters{
	// ...
}


read, err := client.Release(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.RestoreKey`

```go
ctx := context.TODO()

payload := keys.KeyRestoreParameters{
	// ...
}


read, err := client.RestoreKey(ctx, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.RotateKey`

```go
ctx := context.TODO()
id := keys.NewKeyID("keyName")

read, err := client.RotateKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.UnwrapKey`

```go
ctx := context.TODO()
id := keys.NewKeyversionID("https://endpoint-url.example.com", "keyName")

payload := keys.KeyOperationsParameters{
	// ...
}


read, err := client.UnwrapKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.UpdateKey`

```go
ctx := context.TODO()
id := keys.NewKeyversionID("https://endpoint-url.example.com", "keyName")

payload := keys.KeyUpdateParameters{
	// ...
}


read, err := client.UpdateKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.UpdateKeyRotationPolicy`

```go
ctx := context.TODO()
id := keys.NewKeyID("keyName")

payload := keys.KeyRotationPolicy{
	// ...
}


read, err := client.UpdateKeyRotationPolicy(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.Verify`

```go
ctx := context.TODO()
id := keys.NewKeyversionID("https://endpoint-url.example.com", "keyName")

payload := keys.KeyVerifyParameters{
	// ...
}


read, err := client.Verify(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `KeysClient.WrapKey`

```go
ctx := context.TODO()
id := keys.NewKeyversionID("https://endpoint-url.example.com", "keyName")

payload := keys.KeyOperationsParameters{
	// ...
}


read, err := client.WrapKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
