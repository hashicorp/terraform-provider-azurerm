
## `github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/vaults` Documentation

The `vaults` SDK allows for interaction with Azure Resource Manager `keyvault` (API Version `2023-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/vaults"
```


### Client Initialization

```go
client := vaults.NewVaultsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VaultsClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := vaults.VaultCheckNameAvailabilityParameters{
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


### Example Usage: `VaultsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewKeyVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName")

payload := vaults.VaultCreateOrUpdateParameters{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VaultsClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewKeyVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VaultsClient.Get`

```go
ctx := context.TODO()
id := commonids.NewKeyVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VaultsClient.GetDeleted`

```go
ctx := context.TODO()
id := vaults.NewDeletedVaultID("12345678-1234-9876-4563-123456789012", "locationName", "deletedVaultName")

read, err := client.GetDeleted(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VaultsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id, vaults.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, vaults.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VaultsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, vaults.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, vaults.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VaultsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, vaults.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, vaults.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VaultsClient.ListDeleted`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListDeleted(ctx, id)` can be used to do batched pagination
items, err := client.ListDeletedComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VaultsClient.PurgeDeleted`

```go
ctx := context.TODO()
id := vaults.NewDeletedVaultID("12345678-1234-9876-4563-123456789012", "locationName", "deletedVaultName")

if err := client.PurgeDeletedThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VaultsClient.Update`

```go
ctx := context.TODO()
id := commonids.NewKeyVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName")

payload := vaults.VaultPatchParameters{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VaultsClient.UpdateAccessPolicy`

```go
ctx := context.TODO()
id := vaults.NewOperationKindID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "add")

payload := vaults.VaultAccessPolicyParameters{
	// ...
}


read, err := client.UpdateAccessPolicy(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
