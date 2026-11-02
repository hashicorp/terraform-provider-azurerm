
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/disasterrecoveryconfigs` Documentation

The `disasterrecoveryconfigs` SDK allows for interaction with Azure Resource Manager `servicebus` (API Version `2026-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/disasterrecoveryconfigs"
```


### Client Initialization

```go
client := disasterrecoveryconfigs.NewDisasterRecoveryConfigsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DisasterRecoveryConfigsClient.GetAuthorizationRule`

```go
ctx := context.TODO()
id := disasterrecoveryconfigs.NewDisasterRecoveryConfigAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "disasterRecoveryConfigName", "authorizationRuleName")

read, err := client.GetAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DisasterRecoveryConfigsClient.ListAuthorizationRules`

```go
ctx := context.TODO()
id := disasterrecoveryconfigs.NewDisasterRecoveryConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "disasterRecoveryConfigName")

// alternatively `client.ListAuthorizationRules(ctx, id)` can be used to do batched pagination
items, err := client.ListAuthorizationRulesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DisasterRecoveryConfigsClient.ListKeys`

```go
ctx := context.TODO()
id := disasterrecoveryconfigs.NewDisasterRecoveryConfigAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "disasterRecoveryConfigName", "authorizationRuleName")

read, err := client.ListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
