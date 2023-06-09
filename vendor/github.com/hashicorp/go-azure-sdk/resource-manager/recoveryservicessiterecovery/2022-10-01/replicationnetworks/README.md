
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationnetworks` Documentation

The `replicationnetworks` SDK allows for interaction with the Azure Resource Manager Service `recoveryservicessiterecovery` (API Version `2022-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationnetworks"
```


### Client Initialization

```go
client := replicationnetworks.NewReplicationNetworksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReplicationNetworksClient.Get`

```go
ctx := context.TODO()
id := replicationnetworks.NewReplicationNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue", "replicationNetworkValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReplicationNetworksClient.List`

```go
ctx := context.TODO()
id := replicationnetworks.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationNetworksClient.ListByReplicationFabrics`

```go
ctx := context.TODO()
id := replicationnetworks.NewReplicationFabricID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue")

// alternatively `client.ListByReplicationFabrics(ctx, id)` can be used to do batched pagination
items, err := client.ListByReplicationFabricsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
