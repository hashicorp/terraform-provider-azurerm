
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationfabrics` Documentation

The `replicationfabrics` SDK allows for interaction with the Azure Resource Manager Service `recoveryservicessiterecovery` (API Version `2022-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationfabrics"
```


### Client Initialization

```go
client := replicationfabrics.NewReplicationFabricsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReplicationFabricsClient.CheckConsistency`

```go
ctx := context.TODO()
id := replicationfabrics.NewReplicationFabricID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue")

if err := client.CheckConsistencyThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationFabricsClient.Create`

```go
ctx := context.TODO()
id := replicationfabrics.NewReplicationFabricID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue")

payload := replicationfabrics.FabricCreationInput{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationFabricsClient.Delete`

```go
ctx := context.TODO()
id := replicationfabrics.NewReplicationFabricID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationFabricsClient.Get`

```go
ctx := context.TODO()
id := replicationfabrics.NewReplicationFabricID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue")

read, err := client.Get(ctx, id, replicationfabrics.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReplicationFabricsClient.List`

```go
ctx := context.TODO()
id := replicationfabrics.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationFabricsClient.MigrateToAad`

```go
ctx := context.TODO()
id := replicationfabrics.NewReplicationFabricID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue")

if err := client.MigrateToAadThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationFabricsClient.Purge`

```go
ctx := context.TODO()
id := replicationfabrics.NewReplicationFabricID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue")

if err := client.PurgeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationFabricsClient.ReassociateGateway`

```go
ctx := context.TODO()
id := replicationfabrics.NewReplicationFabricID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue")

payload := replicationfabrics.FailoverProcessServerRequest{
	// ...
}


if err := client.ReassociateGatewayThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationFabricsClient.RenewCertificate`

```go
ctx := context.TODO()
id := replicationfabrics.NewReplicationFabricID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationFabricValue")

payload := replicationfabrics.RenewCertificateInput{
	// ...
}


if err := client.RenewCertificateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
