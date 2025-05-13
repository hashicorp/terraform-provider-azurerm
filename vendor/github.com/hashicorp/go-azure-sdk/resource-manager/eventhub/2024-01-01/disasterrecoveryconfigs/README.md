
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/disasterrecoveryconfigs` Documentation

The `disasterrecoveryconfigs` SDK allows for interaction with Azure Resource Manager `eventhub` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/disasterrecoveryconfigs"
```


### Client Initialization

```go
client := disasterrecoveryconfigs.NewDisasterRecoveryConfigsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DisasterRecoveryConfigsClient.BreakPairing`

```go
ctx := context.TODO()
id := disasterrecoveryconfigs.NewDisasterRecoveryConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "disasterRecoveryConfigName")

read, err := client.BreakPairing(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DisasterRecoveryConfigsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := disasterrecoveryconfigs.NewDisasterRecoveryConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "disasterRecoveryConfigName")

payload := disasterrecoveryconfigs.ArmDisasterRecovery{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DisasterRecoveryConfigsClient.Delete`

```go
ctx := context.TODO()
id := disasterrecoveryconfigs.NewDisasterRecoveryConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "disasterRecoveryConfigName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DisasterRecoveryConfigsClient.FailOver`

```go
ctx := context.TODO()
id := disasterrecoveryconfigs.NewDisasterRecoveryConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "disasterRecoveryConfigName")

read, err := client.FailOver(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DisasterRecoveryConfigsClient.Get`

```go
ctx := context.TODO()
id := disasterrecoveryconfigs.NewDisasterRecoveryConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "disasterRecoveryConfigName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DisasterRecoveryConfigsClient.List`

```go
ctx := context.TODO()
id := disasterrecoveryconfigs.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
