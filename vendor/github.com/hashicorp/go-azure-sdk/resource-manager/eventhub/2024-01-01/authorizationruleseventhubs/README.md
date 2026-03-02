
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/authorizationruleseventhubs` Documentation

The `authorizationruleseventhubs` SDK allows for interaction with Azure Resource Manager `eventhub` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/authorizationruleseventhubs"
```


### Client Initialization

```go
client := authorizationruleseventhubs.NewAuthorizationRulesEventHubsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AuthorizationRulesEventHubsClient.EventHubsCreateOrUpdateAuthorizationRule`

```go
ctx := context.TODO()
id := authorizationruleseventhubs.NewEventhubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "eventhubName", "authorizationRuleName")

payload := authorizationruleseventhubs.AuthorizationRule{
	// ...
}


read, err := client.EventHubsCreateOrUpdateAuthorizationRule(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationRulesEventHubsClient.EventHubsListAuthorizationRules`

```go
ctx := context.TODO()
id := authorizationruleseventhubs.NewEventhubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "eventhubName")

// alternatively `client.EventHubsListAuthorizationRules(ctx, id)` can be used to do batched pagination
items, err := client.EventHubsListAuthorizationRulesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AuthorizationRulesEventHubsClient.EventHubsListKeys`

```go
ctx := context.TODO()
id := authorizationruleseventhubs.NewEventhubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "eventhubName", "authorizationRuleName")

read, err := client.EventHubsListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationRulesEventHubsClient.EventHubsRegenerateKeys`

```go
ctx := context.TODO()
id := authorizationruleseventhubs.NewEventhubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "eventhubName", "authorizationRuleName")

payload := authorizationruleseventhubs.RegenerateAccessKeyParameters{
	// ...
}


read, err := client.EventHubsRegenerateKeys(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
