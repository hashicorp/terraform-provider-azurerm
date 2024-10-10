
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainers` Documentation

The `replicationprotectioncontainers` SDK allows for interaction with Azure Resource Manager `recoveryservicessiterecovery` (API Version `2024-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainers"
```


### Client Initialization

```go
client := replicationprotectioncontainers.NewReplicationProtectionContainersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReplicationProtectionContainersClient.Create`

```go
ctx := context.TODO()
id := replicationprotectioncontainers.NewReplicationProtectionContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName")

payload := replicationprotectioncontainers.CreateProtectionContainerInput{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectionContainersClient.Delete`

```go
ctx := context.TODO()
id := replicationprotectioncontainers.NewReplicationProtectionContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectionContainersClient.DiscoverProtectableItem`

```go
ctx := context.TODO()
id := replicationprotectioncontainers.NewReplicationProtectionContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName")

payload := replicationprotectioncontainers.DiscoverProtectableItemRequest{
	// ...
}


if err := client.DiscoverProtectableItemThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectionContainersClient.Get`

```go
ctx := context.TODO()
id := replicationprotectioncontainers.NewReplicationProtectionContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReplicationProtectionContainersClient.List`

```go
ctx := context.TODO()
id := replicationprotectioncontainers.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationProtectionContainersClient.ListByReplicationFabrics`

```go
ctx := context.TODO()
id := replicationprotectioncontainers.NewReplicationFabricID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName")

// alternatively `client.ListByReplicationFabrics(ctx, id)` can be used to do batched pagination
items, err := client.ListByReplicationFabricsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationProtectionContainersClient.SwitchClusterProtection`

```go
ctx := context.TODO()
id := replicationprotectioncontainers.NewReplicationProtectionContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName")

payload := replicationprotectioncontainers.SwitchClusterProtectionInput{
	// ...
}


if err := client.SwitchClusterProtectionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectionContainersClient.SwitchProtection`

```go
ctx := context.TODO()
id := replicationprotectioncontainers.NewReplicationProtectionContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName")

payload := replicationprotectioncontainers.SwitchProtectionInput{
	// ...
}


if err := client.SwitchProtectionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
