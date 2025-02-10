
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2024-01-01/vaults` Documentation

The `vaults` SDK allows for interaction with Azure Resource Manager `recoveryservices` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2024-01-01/vaults"
```


### Client Initialization

```go
client := vaults.NewVaultsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VaultsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := vaults.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName")

payload := vaults.Vault{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VaultsClient.Delete`

```go
ctx := context.TODO()
id := vaults.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName")

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
id := vaults.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VaultsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VaultsClient.ListBySubscriptionId`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscriptionId(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionIdComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VaultsClient.Update`

```go
ctx := context.TODO()
id := vaults.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName")

payload := vaults.PatchVault{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
