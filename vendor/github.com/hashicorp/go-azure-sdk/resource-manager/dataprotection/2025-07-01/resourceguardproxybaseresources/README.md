
## `github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardproxybaseresources` Documentation

The `resourceguardproxybaseresources` SDK allows for interaction with Azure Resource Manager `dataprotection` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardproxybaseresources"
```


### Client Initialization

```go
client := resourceguardproxybaseresources.NewResourceGuardProxyBaseResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ResourceGuardProxyBaseResourcesClient.DppResourceGuardProxyCreateOrUpdate`

```go
ctx := context.TODO()
id := resourceguardproxybaseresources.NewBackupResourceGuardProxyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupResourceGuardProxyName")

payload := resourceguardproxybaseresources.ResourceGuardProxyBaseResource{
	// ...
}


read, err := client.DppResourceGuardProxyCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceGuardProxyBaseResourcesClient.DppResourceGuardProxyDelete`

```go
ctx := context.TODO()
id := resourceguardproxybaseresources.NewBackupResourceGuardProxyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupResourceGuardProxyName")

read, err := client.DppResourceGuardProxyDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceGuardProxyBaseResourcesClient.DppResourceGuardProxyGet`

```go
ctx := context.TODO()
id := resourceguardproxybaseresources.NewBackupResourceGuardProxyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupResourceGuardProxyName")

read, err := client.DppResourceGuardProxyGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceGuardProxyBaseResourcesClient.DppResourceGuardProxyList`

```go
ctx := context.TODO()
id := resourceguardproxybaseresources.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName")

// alternatively `client.DppResourceGuardProxyList(ctx, id)` can be used to do batched pagination
items, err := client.DppResourceGuardProxyListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceGuardProxyBaseResourcesClient.DppResourceGuardProxyUnlockDelete`

```go
ctx := context.TODO()
id := resourceguardproxybaseresources.NewBackupResourceGuardProxyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupResourceGuardProxyName")

payload := resourceguardproxybaseresources.UnlockDeleteRequest{
	// ...
}


read, err := client.DppResourceGuardProxyUnlockDelete(ctx, id, payload, resourceguardproxybaseresources.DefaultDppResourceGuardProxyUnlockDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
