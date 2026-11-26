
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/armdisasterrecoveries` Documentation

The `armdisasterrecoveries` SDK allows for interaction with Azure Resource Manager `servicebus` (API Version `2026-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/armdisasterrecoveries"
```


### Client Initialization

```go
client := armdisasterrecoveries.NewArmDisasterRecoveriesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ArmDisasterRecoveriesClient.DisasterRecoveryConfigsBreakPairing`

```go
ctx := context.TODO()
id := armdisasterrecoveries.NewDisasterRecoveryConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "disasterRecoveryConfigName")

read, err := client.DisasterRecoveryConfigsBreakPairing(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ArmDisasterRecoveriesClient.DisasterRecoveryConfigsCreateOrUpdate`

```go
ctx := context.TODO()
id := armdisasterrecoveries.NewDisasterRecoveryConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "disasterRecoveryConfigName")

payload := armdisasterrecoveries.ArmDisasterRecovery{
	// ...
}


read, err := client.DisasterRecoveryConfigsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ArmDisasterRecoveriesClient.DisasterRecoveryConfigsDelete`

```go
ctx := context.TODO()
id := armdisasterrecoveries.NewDisasterRecoveryConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "disasterRecoveryConfigName")

read, err := client.DisasterRecoveryConfigsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ArmDisasterRecoveriesClient.DisasterRecoveryConfigsFailOver`

```go
ctx := context.TODO()
id := armdisasterrecoveries.NewDisasterRecoveryConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "disasterRecoveryConfigName")

payload := armdisasterrecoveries.NamespaceFailoverProperties{
	// ...
}


read, err := client.DisasterRecoveryConfigsFailOver(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ArmDisasterRecoveriesClient.DisasterRecoveryConfigsGet`

```go
ctx := context.TODO()
id := armdisasterrecoveries.NewDisasterRecoveryConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "disasterRecoveryConfigName")

read, err := client.DisasterRecoveryConfigsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ArmDisasterRecoveriesClient.DisasterRecoveryConfigsList`

```go
ctx := context.TODO()
id := armdisasterrecoveries.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

// alternatively `client.DisasterRecoveryConfigsList(ctx, id)` can be used to do batched pagination
items, err := client.DisasterRecoveryConfigsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
