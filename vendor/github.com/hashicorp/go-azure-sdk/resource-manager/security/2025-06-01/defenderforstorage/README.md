
## `github.com/hashicorp/go-azure-sdk/resource-manager/security/2025-06-01/defenderforstorage` Documentation

The `defenderforstorage` SDK allows for interaction with Azure Resource Manager `security` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/security/2025-06-01/defenderforstorage"
```


### Client Initialization

```go
client := defenderforstorage.NewDefenderForStorageClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DefenderForStorageClient.CancelMalwareScan`

```go
ctx := context.TODO()
id := defenderforstorage.NewScopedMalwareScanID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "scanId")

read, err := client.CancelMalwareScan(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DefenderForStorageClient.Create`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := defenderforstorage.DefenderForStorageSetting{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DefenderForStorageClient.Get`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DefenderForStorageClient.GetMalwareScan`

```go
ctx := context.TODO()
id := defenderforstorage.NewScopedMalwareScanID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "scanId")

read, err := client.GetMalwareScan(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DefenderForStorageClient.StartMalwareScan`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.StartMalwareScan(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
