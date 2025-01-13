
## `github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/replicas` Documentation

The `replicas` SDK allows for interaction with Azure Resource Manager `appconfiguration` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/replicas"
```


### Client Initialization

```go
client := replicas.NewReplicasClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReplicasClient.Create`

```go
ctx := context.TODO()
id := replicas.NewReplicaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configurationStoreName", "replicaName")

payload := replicas.Replica{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicasClient.Delete`

```go
ctx := context.TODO()
id := replicas.NewReplicaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configurationStoreName", "replicaName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicasClient.Get`

```go
ctx := context.TODO()
id := replicas.NewReplicaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configurationStoreName", "replicaName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReplicasClient.ListByConfigurationStore`

```go
ctx := context.TODO()
id := replicas.NewConfigurationStoreID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configurationStoreName")

// alternatively `client.ListByConfigurationStore(ctx, id)` can be used to do batched pagination
items, err := client.ListByConfigurationStoreComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
