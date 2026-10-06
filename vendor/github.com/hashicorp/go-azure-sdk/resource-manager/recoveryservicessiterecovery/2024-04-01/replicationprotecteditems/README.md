
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems` Documentation

The `replicationprotecteditems` SDK allows for interaction with Azure Resource Manager `recoveryservicessiterecovery` (API Version `2024-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
```


### Client Initialization

```go
client := replicationprotecteditems.NewReplicationProtectedItemsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReplicationProtectedItemsClient.AddDisks`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

payload := replicationprotecteditems.AddDisksInput{
	// ...
}


if err := client.AddDisksThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.ApplyRecoveryPoint`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

payload := replicationprotecteditems.ApplyRecoveryPointInput{
	// ...
}


if err := client.ApplyRecoveryPointThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.Create`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

payload := replicationprotecteditems.EnableProtectionInput{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.Delete`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

payload := replicationprotecteditems.DisableProtectionInput{
	// ...
}


if err := client.DeleteThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.FailoverCancel`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

if err := client.FailoverCancelThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.FailoverCommit`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

if err := client.FailoverCommitThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.Get`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReplicationProtectedItemsClient.List`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName")

// alternatively `client.List(ctx, id, replicationprotecteditems.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, replicationprotecteditems.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationProtectedItemsClient.ListByReplicationProtectionContainers`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectionContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName")

// alternatively `client.ListByReplicationProtectionContainers(ctx, id)` can be used to do batched pagination
items, err := client.ListByReplicationProtectionContainersComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationProtectedItemsClient.PlannedFailover`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

payload := replicationprotecteditems.PlannedFailoverInput{
	// ...
}


if err := client.PlannedFailoverThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.Purge`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

if err := client.PurgeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.RemoveDisks`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

payload := replicationprotecteditems.RemoveDisksInput{
	// ...
}


if err := client.RemoveDisksThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.RepairReplication`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

if err := client.RepairReplicationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.Reprotect`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

payload := replicationprotecteditems.ReverseReplicationInput{
	// ...
}


if err := client.ReprotectThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.ResolveHealthErrors`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

payload := replicationprotecteditems.ResolveHealthInput{
	// ...
}


if err := client.ResolveHealthErrorsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.SwitchProvider`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

payload := replicationprotecteditems.SwitchProviderInput{
	// ...
}


if err := client.SwitchProviderThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.TestFailover`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

payload := replicationprotecteditems.TestFailoverInput{
	// ...
}


if err := client.TestFailoverThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.TestFailoverCleanup`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

payload := replicationprotecteditems.TestFailoverCleanupInput{
	// ...
}


if err := client.TestFailoverCleanupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.UnplannedFailover`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

payload := replicationprotecteditems.UnplannedFailoverInput{
	// ...
}


if err := client.UnplannedFailoverThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.Update`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

payload := replicationprotecteditems.UpdateReplicationProtectedItemInput{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.UpdateAppliance`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

payload := replicationprotecteditems.UpdateApplianceForReplicationProtectedItemInput{
	// ...
}


if err := client.UpdateApplianceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationProtectedItemsClient.UpdateMobilityService`

```go
ctx := context.TODO()
id := replicationprotecteditems.NewReplicationProtectedItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "replicationFabricName", "replicationProtectionContainerName", "replicationProtectedItemName")

payload := replicationprotecteditems.UpdateMobilityServiceRequest{
	// ...
}


if err := client.UpdateMobilityServiceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
