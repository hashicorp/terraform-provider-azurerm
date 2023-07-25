
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationnetworkmappings` Documentation

The `replicationnetworkmappings` SDK allows for interaction with the Azure Resource Manager Service `recoveryservicessiterecovery` (API Version `2022-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationnetworkmappings"
```


### Client Initialization

```go
client := replicationnetworkmappings.NewReplicationNetworkMappingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReplicationNetworkMappingsClient.Create`

```go
ctx := context.TODO()
id := replicationnetworkmappings.NewReplicationNetworkMappingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue", "replicationNetworkValue", "replicationNetworkMappingValue")

payload := replicationnetworkmappings.CreateNetworkMappingInput{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationNetworkMappingsClient.Delete`

```go
ctx := context.TODO()
id := replicationnetworkmappings.NewReplicationNetworkMappingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue", "replicationNetworkValue", "replicationNetworkMappingValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationNetworkMappingsClient.Get`

```go
ctx := context.TODO()
id := replicationnetworkmappings.NewReplicationNetworkMappingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue", "replicationNetworkValue", "replicationNetworkMappingValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReplicationNetworkMappingsClient.List`

```go
ctx := context.TODO()
id := replicationnetworkmappings.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationNetworkMappingsClient.ListByReplicationNetworks`

```go
ctx := context.TODO()
id := replicationnetworkmappings.NewReplicationNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue", "replicationNetworkValue")

// alternatively `client.ListByReplicationNetworks(ctx, id)` can be used to do batched pagination
items, err := client.ListByReplicationNetworksComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationNetworkMappingsClient.Update`

```go
ctx := context.TODO()
id := replicationnetworkmappings.NewReplicationNetworkMappingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue", "replicationNetworkValue", "replicationNetworkMappingValue")

payload := replicationnetworkmappings.UpdateNetworkMappingInput{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
