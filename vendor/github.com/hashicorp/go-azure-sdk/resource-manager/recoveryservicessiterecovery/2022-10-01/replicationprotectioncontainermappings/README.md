
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectioncontainermappings` Documentation

The `replicationprotectioncontainermappings` SDK allows for interaction with the Azure Resource Manager Service `recoveryservicessiterecovery` (API Version `2022-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectioncontainermappings"
```


### Client Initialization

```go
client := replicationprotectioncontainermappings.NewReplicationProtectionContainerMappingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReplicationProtectionContainerMappingsClient.Create`

```go
ctx := context.TODO()
id := replicationprotectioncontainermappings.NewReplicationProtectionContainerMappingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue", "replicationProtectionContainerValue", "replicationProtectionContainerMappingValue")

payload := replicationprotectioncontainermappings.CreateProtectionContainerMappingInput{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectionContainerMappingsClient.Delete`

```go
ctx := context.TODO()
id := replicationprotectioncontainermappings.NewReplicationProtectionContainerMappingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue", "replicationProtectionContainerValue", "replicationProtectionContainerMappingValue")

payload := replicationprotectioncontainermappings.RemoveProtectionContainerMappingInput{
	// ...
}


if err := client.DeleteThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectionContainerMappingsClient.Get`

```go
ctx := context.TODO()
id := replicationprotectioncontainermappings.NewReplicationProtectionContainerMappingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue", "replicationProtectionContainerValue", "replicationProtectionContainerMappingValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReplicationProtectionContainerMappingsClient.List`

```go
ctx := context.TODO()
id := replicationprotectioncontainermappings.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationProtectionContainerMappingsClient.ListByReplicationProtectionContainers`

```go
ctx := context.TODO()
id := replicationprotectioncontainermappings.NewReplicationProtectionContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue", "replicationProtectionContainerValue")

// alternatively `client.ListByReplicationProtectionContainers(ctx, id)` can be used to do batched pagination
items, err := client.ListByReplicationProtectionContainersComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationProtectionContainerMappingsClient.Purge`

```go
ctx := context.TODO()
id := replicationprotectioncontainermappings.NewReplicationProtectionContainerMappingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue", "replicationProtectionContainerValue", "replicationProtectionContainerMappingValue")

if err := client.PurgeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectionContainerMappingsClient.Update`

```go
ctx := context.TODO()
id := replicationprotectioncontainermappings.NewReplicationProtectionContainerMappingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue", "replicationProtectionContainerValue", "replicationProtectionContainerMappingValue")

payload := replicationprotectioncontainermappings.UpdateProtectionContainerMappingInput{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
