
## `github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/deletedmanagedhsms` Documentation

The `deletedmanagedhsms` SDK allows for interaction with Azure Resource Manager `keyvault` (API Version `2026-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/deletedmanagedhsms"
```


### Client Initialization

```go
client := deletedmanagedhsms.NewDeletedManagedHsmsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeletedManagedHsmsClient.ManagedHsmsGetDeleted`

```go
ctx := context.TODO()
id := deletedmanagedhsms.NewDeletedManagedHSMID("12345678-1234-9876-4563-123456789012", "locationName", "deletedManagedHSMName")

read, err := client.ManagedHsmsGetDeleted(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeletedManagedHsmsClient.ManagedHsmsPurgeDeleted`

```go
ctx := context.TODO()
id := deletedmanagedhsms.NewDeletedManagedHSMID("12345678-1234-9876-4563-123456789012", "locationName", "deletedManagedHSMName")

if err := client.ManagedHsmsPurgeDeletedThenPoll(ctx, id); err != nil {
	// handle the error
}
```
