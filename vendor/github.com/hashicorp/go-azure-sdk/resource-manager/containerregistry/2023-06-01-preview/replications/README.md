
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/replications` Documentation

The `replications` SDK allows for interaction with Azure Resource Manager `containerregistry` (API Version `2023-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/replications"
```


### Client Initialization

```go
client := replications.NewReplicationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReplicationsClient.Create`

```go
ctx := context.TODO()
id := replications.NewReplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "replicationName")

payload := replications.Replication{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationsClient.Delete`

```go
ctx := context.TODO()
id := replications.NewReplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "replicationName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationsClient.Get`

```go
ctx := context.TODO()
id := replications.NewReplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "replicationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReplicationsClient.List`

```go
ctx := context.TODO()
id := replications.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationsClient.Update`

```go
ctx := context.TODO()
id := replications.NewReplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "replicationName")

payload := replications.ReplicationUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
