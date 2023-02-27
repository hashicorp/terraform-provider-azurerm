
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectableitems` Documentation

The `replicationprotectableitems` SDK allows for interaction with the Azure Resource Manager Service `recoveryservicessiterecovery` (API Version `2022-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectableitems"
```


### Client Initialization

```go
client := replicationprotectableitems.NewReplicationProtectableItemsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReplicationProtectableItemsClient.Get`

```go
ctx := context.TODO()
id := replicationprotectableitems.NewReplicationProtectableItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue", "replicationProtectionContainerValue", "replicationProtectableItemValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReplicationProtectableItemsClient.ListByReplicationProtectionContainers`

```go
ctx := context.TODO()
id := replicationprotectableitems.NewReplicationProtectionContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue", "replicationProtectionContainerValue")

// alternatively `client.ListByReplicationProtectionContainers(ctx, id, replicationprotectableitems.DefaultListByReplicationProtectionContainersOperationOptions())` can be used to do batched pagination
items, err := client.ListByReplicationProtectionContainersComplete(ctx, id, replicationprotectableitems.DefaultListByReplicationProtectionContainersOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
