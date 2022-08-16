
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queuesauthorizationrule` Documentation

The `queuesauthorizationrule` SDK allows for interaction with the Azure Resource Manager Service `servicebus` (API Version `2021-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queuesauthorizationrule"
```


### Client Initialization

```go
client := queuesauthorizationrule.NewQueuesAuthorizationRuleClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `QueuesAuthorizationRuleClient.QueuesCreateOrUpdateAuthorizationRule`

```go
ctx := context.TODO()
id := queuesauthorizationrule.NewQueueAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "queueValue", "authorizationRuleValue")

payload := queuesauthorizationrule.SBAuthorizationRule{
	// ...
}


read, err := client.QueuesCreateOrUpdateAuthorizationRule(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueuesAuthorizationRuleClient.QueuesDeleteAuthorizationRule`

```go
ctx := context.TODO()
id := queuesauthorizationrule.NewQueueAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "queueValue", "authorizationRuleValue")

read, err := client.QueuesDeleteAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueuesAuthorizationRuleClient.QueuesGetAuthorizationRule`

```go
ctx := context.TODO()
id := queuesauthorizationrule.NewQueueAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "queueValue", "authorizationRuleValue")

read, err := client.QueuesGetAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueuesAuthorizationRuleClient.QueuesListAuthorizationRules`

```go
ctx := context.TODO()
id := queuesauthorizationrule.NewQueueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "queueValue")

// alternatively `client.QueuesListAuthorizationRules(ctx, id)` can be used to do batched pagination
items, err := client.QueuesListAuthorizationRulesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `QueuesAuthorizationRuleClient.QueuesListKeys`

```go
ctx := context.TODO()
id := queuesauthorizationrule.NewQueueAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "queueValue", "authorizationRuleValue")

read, err := client.QueuesListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueuesAuthorizationRuleClient.QueuesRegenerateKeys`

```go
ctx := context.TODO()
id := queuesauthorizationrule.NewQueueAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "queueValue", "authorizationRuleValue")

payload := queuesauthorizationrule.RegenerateAccessKeyParameters{
	// ...
}


read, err := client.QueuesRegenerateKeys(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
