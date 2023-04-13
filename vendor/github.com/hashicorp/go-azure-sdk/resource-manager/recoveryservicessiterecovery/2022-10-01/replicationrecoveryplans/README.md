
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationrecoveryplans` Documentation

The `replicationrecoveryplans` SDK allows for interaction with the Azure Resource Manager Service `recoveryservicessiterecovery` (API Version `2022-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationrecoveryplans"
```


### Client Initialization

```go
client := replicationrecoveryplans.NewReplicationRecoveryPlansClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReplicationRecoveryPlansClient.Create`

```go
ctx := context.TODO()
id := replicationrecoveryplans.NewReplicationRecoveryPlanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationRecoveryPlanValue")

payload := replicationrecoveryplans.CreateRecoveryPlanInput{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationRecoveryPlansClient.Delete`

```go
ctx := context.TODO()
id := replicationrecoveryplans.NewReplicationRecoveryPlanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationRecoveryPlanValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationRecoveryPlansClient.FailoverCancel`

```go
ctx := context.TODO()
id := replicationrecoveryplans.NewReplicationRecoveryPlanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationRecoveryPlanValue")

if err := client.FailoverCancelThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationRecoveryPlansClient.FailoverCommit`

```go
ctx := context.TODO()
id := replicationrecoveryplans.NewReplicationRecoveryPlanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationRecoveryPlanValue")

if err := client.FailoverCommitThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationRecoveryPlansClient.Get`

```go
ctx := context.TODO()
id := replicationrecoveryplans.NewReplicationRecoveryPlanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationRecoveryPlanValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReplicationRecoveryPlansClient.List`

```go
ctx := context.TODO()
id := replicationrecoveryplans.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationRecoveryPlansClient.PlannedFailover`

```go
ctx := context.TODO()
id := replicationrecoveryplans.NewReplicationRecoveryPlanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationRecoveryPlanValue")

payload := replicationrecoveryplans.RecoveryPlanPlannedFailoverInput{
	// ...
}


if err := client.PlannedFailoverThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationRecoveryPlansClient.Reprotect`

```go
ctx := context.TODO()
id := replicationrecoveryplans.NewReplicationRecoveryPlanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationRecoveryPlanValue")

if err := client.ReprotectThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationRecoveryPlansClient.TestFailover`

```go
ctx := context.TODO()
id := replicationrecoveryplans.NewReplicationRecoveryPlanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationRecoveryPlanValue")

payload := replicationrecoveryplans.RecoveryPlanTestFailoverInput{
	// ...
}


if err := client.TestFailoverThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationRecoveryPlansClient.TestFailoverCleanup`

```go
ctx := context.TODO()
id := replicationrecoveryplans.NewReplicationRecoveryPlanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationRecoveryPlanValue")

payload := replicationrecoveryplans.RecoveryPlanTestFailoverCleanupInput{
	// ...
}


if err := client.TestFailoverCleanupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationRecoveryPlansClient.UnplannedFailover`

```go
ctx := context.TODO()
id := replicationrecoveryplans.NewReplicationRecoveryPlanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationRecoveryPlanValue")

payload := replicationrecoveryplans.RecoveryPlanUnplannedFailoverInput{
	// ...
}


if err := client.UnplannedFailoverThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationRecoveryPlansClient.Update`

```go
ctx := context.TODO()
id := replicationrecoveryplans.NewReplicationRecoveryPlanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationRecoveryPlanValue")

payload := replicationrecoveryplans.UpdateRecoveryPlanInput{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
