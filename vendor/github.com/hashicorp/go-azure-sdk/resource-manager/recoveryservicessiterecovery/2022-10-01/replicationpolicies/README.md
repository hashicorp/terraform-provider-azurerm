
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationpolicies` Documentation

The `replicationpolicies` SDK allows for interaction with the Azure Resource Manager Service `recoveryservicessiterecovery` (API Version `2022-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationpolicies"
```


### Client Initialization

```go
client := replicationpolicies.NewReplicationPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReplicationPoliciesClient.Create`

```go
ctx := context.TODO()
id := replicationpolicies.NewReplicationPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationPolicyValue")

payload := replicationpolicies.CreatePolicyInput{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationPoliciesClient.Delete`

```go
ctx := context.TODO()
id := replicationpolicies.NewReplicationPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationPolicyValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationPoliciesClient.Get`

```go
ctx := context.TODO()
id := replicationpolicies.NewReplicationPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationPolicyValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReplicationPoliciesClient.List`

```go
ctx := context.TODO()
id := replicationpolicies.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationPoliciesClient.Update`

```go
ctx := context.TODO()
id := replicationpolicies.NewReplicationPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "replicationPolicyValue")

payload := replicationpolicies.UpdatePolicyInput{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
