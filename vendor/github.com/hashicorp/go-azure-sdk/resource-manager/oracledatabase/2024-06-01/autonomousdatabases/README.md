
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases` Documentation

The `autonomousdatabases` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2024-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases"
```


### Client Initialization

```go
client := autonomousdatabases.NewAutonomousDatabasesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AutonomousDatabasesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := autonomousdatabases.NewAutonomousDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autonomousDatabaseName")

payload := autonomousdatabases.AutonomousDatabase{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AutonomousDatabasesClient.Delete`

```go
ctx := context.TODO()
id := autonomousdatabases.NewAutonomousDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autonomousDatabaseName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AutonomousDatabasesClient.Failover`

```go
ctx := context.TODO()
id := autonomousdatabases.NewAutonomousDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autonomousDatabaseName")

payload := autonomousdatabases.PeerDbDetails{
	// ...
}


if err := client.FailoverThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AutonomousDatabasesClient.GenerateWallet`

```go
ctx := context.TODO()
id := autonomousdatabases.NewAutonomousDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autonomousDatabaseName")

payload := autonomousdatabases.GenerateAutonomousDatabaseWalletDetails{
	// ...
}


read, err := client.GenerateWallet(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AutonomousDatabasesClient.Get`

```go
ctx := context.TODO()
id := autonomousdatabases.NewAutonomousDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autonomousDatabaseName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AutonomousDatabasesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AutonomousDatabasesClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AutonomousDatabasesClient.Restore`

```go
ctx := context.TODO()
id := autonomousdatabases.NewAutonomousDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autonomousDatabaseName")

payload := autonomousdatabases.RestoreAutonomousDatabaseDetails{
	// ...
}


if err := client.RestoreThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AutonomousDatabasesClient.Shrink`

```go
ctx := context.TODO()
id := autonomousdatabases.NewAutonomousDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autonomousDatabaseName")

if err := client.ShrinkThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AutonomousDatabasesClient.Switchover`

```go
ctx := context.TODO()
id := autonomousdatabases.NewAutonomousDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autonomousDatabaseName")

payload := autonomousdatabases.PeerDbDetails{
	// ...
}


if err := client.SwitchoverThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AutonomousDatabasesClient.Update`

```go
ctx := context.TODO()
id := autonomousdatabases.NewAutonomousDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autonomousDatabaseName")

payload := autonomousdatabases.AutonomousDatabaseUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
