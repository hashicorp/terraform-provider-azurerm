
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationrecoveryservicesproviders` Documentation

The `replicationrecoveryservicesproviders` SDK allows for interaction with Azure Resource Manager `recoveryservicessiterecovery` (API Version `2024-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationrecoveryservicesproviders"
```


### Client Initialization

```go
client := replicationrecoveryservicesproviders.NewReplicationRecoveryServicesProvidersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReplicationRecoveryServicesProvidersClient.Create`

```go
ctx := context.TODO()
id := replicationrecoveryservicesproviders.NewReplicationRecoveryServicesProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationRecoveryServicesProviderName")

payload := replicationrecoveryservicesproviders.AddRecoveryServicesProviderInput{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationRecoveryServicesProvidersClient.Delete`

```go
ctx := context.TODO()
id := replicationrecoveryservicesproviders.NewReplicationRecoveryServicesProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationRecoveryServicesProviderName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationRecoveryServicesProvidersClient.Get`

```go
ctx := context.TODO()
id := replicationrecoveryservicesproviders.NewReplicationRecoveryServicesProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationRecoveryServicesProviderName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReplicationRecoveryServicesProvidersClient.List`

```go
ctx := context.TODO()
id := replicationrecoveryservicesproviders.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationRecoveryServicesProvidersClient.ListByReplicationFabrics`

```go
ctx := context.TODO()
id := replicationrecoveryservicesproviders.NewReplicationFabricID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName")

// alternatively `client.ListByReplicationFabrics(ctx, id)` can be used to do batched pagination
items, err := client.ListByReplicationFabricsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationRecoveryServicesProvidersClient.Purge`

```go
ctx := context.TODO()
id := replicationrecoveryservicesproviders.NewReplicationRecoveryServicesProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationRecoveryServicesProviderName")

if err := client.PurgeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationRecoveryServicesProvidersClient.RefreshProvider`

```go
ctx := context.TODO()
id := replicationrecoveryservicesproviders.NewReplicationRecoveryServicesProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationRecoveryServicesProviderName")

if err := client.RefreshProviderThenPoll(ctx, id); err != nil {
	// handle the error
}
```
