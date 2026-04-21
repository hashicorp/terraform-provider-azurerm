
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/objectreplicationpolicyoperationgroup` Documentation

The `objectreplicationpolicyoperationgroup` SDK allows for interaction with Azure Resource Manager `storage` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/objectreplicationpolicyoperationgroup"
```


### Client Initialization

```go
client := objectreplicationpolicyoperationgroup.NewObjectReplicationPolicyOperationGroupClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ObjectReplicationPolicyOperationGroupClient.ObjectReplicationPoliciesCreateOrUpdate`

```go
ctx := context.TODO()
id := objectreplicationpolicyoperationgroup.NewObjectReplicationPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "objectReplicationPolicyId")

payload := objectreplicationpolicyoperationgroup.ObjectReplicationPolicy{
	// ...
}


read, err := client.ObjectReplicationPoliciesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ObjectReplicationPolicyOperationGroupClient.ObjectReplicationPoliciesDelete`

```go
ctx := context.TODO()
id := objectreplicationpolicyoperationgroup.NewObjectReplicationPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "objectReplicationPolicyId")

read, err := client.ObjectReplicationPoliciesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ObjectReplicationPolicyOperationGroupClient.ObjectReplicationPoliciesGet`

```go
ctx := context.TODO()
id := objectreplicationpolicyoperationgroup.NewObjectReplicationPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "objectReplicationPolicyId")

read, err := client.ObjectReplicationPoliciesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ObjectReplicationPolicyOperationGroupClient.ObjectReplicationPoliciesList`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

// alternatively `client.ObjectReplicationPoliciesList(ctx, id)` can be used to do batched pagination
items, err := client.ObjectReplicationPoliciesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
