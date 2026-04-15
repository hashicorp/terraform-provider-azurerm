
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/localuseroperationgroup` Documentation

The `localuseroperationgroup` SDK allows for interaction with Azure Resource Manager `storage` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/localuseroperationgroup"
```


### Client Initialization

```go
client := localuseroperationgroup.NewLocalUserOperationGroupClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LocalUserOperationGroupClient.LocalUsersCreateOrUpdate`

```go
ctx := context.TODO()
id := localuseroperationgroup.NewLocalUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "localUserName")

payload := localuseroperationgroup.LocalUser{
	// ...
}


read, err := client.LocalUsersCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalUserOperationGroupClient.LocalUsersDelete`

```go
ctx := context.TODO()
id := localuseroperationgroup.NewLocalUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "localUserName")

read, err := client.LocalUsersDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalUserOperationGroupClient.LocalUsersGet`

```go
ctx := context.TODO()
id := localuseroperationgroup.NewLocalUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "localUserName")

read, err := client.LocalUsersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalUserOperationGroupClient.LocalUsersList`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

// alternatively `client.LocalUsersList(ctx, id, localuseroperationgroup.DefaultLocalUsersListOperationOptions())` can be used to do batched pagination
items, err := client.LocalUsersListComplete(ctx, id, localuseroperationgroup.DefaultLocalUsersListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalUserOperationGroupClient.LocalUsersListKeys`

```go
ctx := context.TODO()
id := localuseroperationgroup.NewLocalUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "localUserName")

read, err := client.LocalUsersListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalUserOperationGroupClient.LocalUsersRegeneratePassword`

```go
ctx := context.TODO()
id := localuseroperationgroup.NewLocalUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "localUserName")

read, err := client.LocalUsersRegeneratePassword(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
