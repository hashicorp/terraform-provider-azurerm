
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/immutabilitypolicies` Documentation

The `immutabilitypolicies` SDK allows for interaction with Azure Resource Manager `storage` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/immutabilitypolicies"
```


### Client Initialization

```go
client := immutabilitypolicies.NewImmutabilityPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ImmutabilityPoliciesClient.BlobContainersCreateOrUpdateImmutabilityPolicy`

```go
ctx := context.TODO()
id := commonids.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "containerName")

payload := immutabilitypolicies.ImmutabilityPolicy{
	// ...
}


read, err := client.BlobContainersCreateOrUpdateImmutabilityPolicy(ctx, id, payload, immutabilitypolicies.DefaultBlobContainersCreateOrUpdateImmutabilityPolicyOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ImmutabilityPoliciesClient.BlobContainersDeleteImmutabilityPolicy`

```go
ctx := context.TODO()
id := commonids.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "containerName")

read, err := client.BlobContainersDeleteImmutabilityPolicy(ctx, id, immutabilitypolicies.DefaultBlobContainersDeleteImmutabilityPolicyOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ImmutabilityPoliciesClient.BlobContainersExtendImmutabilityPolicy`

```go
ctx := context.TODO()
id := commonids.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "containerName")

payload := immutabilitypolicies.ImmutabilityPolicy{
	// ...
}


read, err := client.BlobContainersExtendImmutabilityPolicy(ctx, id, payload, immutabilitypolicies.DefaultBlobContainersExtendImmutabilityPolicyOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ImmutabilityPoliciesClient.BlobContainersGetImmutabilityPolicy`

```go
ctx := context.TODO()
id := commonids.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "containerName")

read, err := client.BlobContainersGetImmutabilityPolicy(ctx, id, immutabilitypolicies.DefaultBlobContainersGetImmutabilityPolicyOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ImmutabilityPoliciesClient.BlobContainersLockImmutabilityPolicy`

```go
ctx := context.TODO()
id := commonids.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "containerName")

read, err := client.BlobContainersLockImmutabilityPolicy(ctx, id, immutabilitypolicies.DefaultBlobContainersLockImmutabilityPolicyOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
