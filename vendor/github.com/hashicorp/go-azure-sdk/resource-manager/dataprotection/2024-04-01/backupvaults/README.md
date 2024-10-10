
## `github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupvaults` Documentation

The `backupvaults` SDK allows for interaction with Azure Resource Manager `dataprotection` (API Version `2024-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupvaults"
```


### Client Initialization

```go
client := backupvaults.NewBackupVaultsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackupVaultsClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := backupvaults.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName")

payload := backupvaults.CheckNameAvailabilityRequest{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupVaultsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := backupvaults.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName")

payload := backupvaults.BackupVaultResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, backupvaults.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `BackupVaultsClient.Delete`

```go
ctx := context.TODO()
id := backupvaults.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupVaultsClient.Get`

```go
ctx := context.TODO()
id := backupvaults.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupVaultsClient.GetInResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.GetInResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.GetInResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BackupVaultsClient.GetInSubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.GetInSubscription(ctx, id)` can be used to do batched pagination
items, err := client.GetInSubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BackupVaultsClient.Update`

```go
ctx := context.TODO()
id := backupvaults.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName")

payload := backupvaults.PatchResourceRequestInput{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload, backupvaults.DefaultUpdateOperationOptions()); err != nil {
	// handle the error
}
```
