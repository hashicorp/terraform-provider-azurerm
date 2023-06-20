
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationvaultsetting` Documentation

The `replicationvaultsetting` SDK allows for interaction with the Azure Resource Manager Service `recoveryservicessiterecovery` (API Version `2022-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationvaultsetting"
```


### Client Initialization

```go
client := replicationvaultsetting.NewReplicationVaultSettingClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReplicationVaultSettingClient.Create`

```go
ctx := context.TODO()
id := replicationvaultsetting.NewReplicationVaultSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationVaultSettingValue")

payload := replicationvaultsetting.VaultSettingCreationInput{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationVaultSettingClient.Get`

```go
ctx := context.TODO()
id := replicationvaultsetting.NewReplicationVaultSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationVaultSettingValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReplicationVaultSettingClient.List`

```go
ctx := context.TODO()
id := replicationvaultsetting.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
