
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/topicsauthorizationrule` Documentation

The `topicsauthorizationrule` SDK allows for interaction with Azure Resource Manager `servicebus` (API Version `2021-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/topicsauthorizationrule"
```


### Client Initialization

```go
client := topicsauthorizationrule.NewTopicsAuthorizationRuleClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TopicsAuthorizationRuleClient.TopicsCreateOrUpdateAuthorizationRule`

```go
ctx := context.TODO()
id := topicsauthorizationrule.NewTopicAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "topicName", "authorizationRuleName")

payload := topicsauthorizationrule.SBAuthorizationRule{
	// ...
}


read, err := client.TopicsCreateOrUpdateAuthorizationRule(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TopicsAuthorizationRuleClient.TopicsDeleteAuthorizationRule`

```go
ctx := context.TODO()
id := topicsauthorizationrule.NewTopicAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "topicName", "authorizationRuleName")

read, err := client.TopicsDeleteAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TopicsAuthorizationRuleClient.TopicsGetAuthorizationRule`

```go
ctx := context.TODO()
id := topicsauthorizationrule.NewTopicAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "topicName", "authorizationRuleName")

read, err := client.TopicsGetAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TopicsAuthorizationRuleClient.TopicsListAuthorizationRules`

```go
ctx := context.TODO()
id := topicsauthorizationrule.NewTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "topicName")

// alternatively `client.TopicsListAuthorizationRules(ctx, id)` can be used to do batched pagination
items, err := client.TopicsListAuthorizationRulesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TopicsAuthorizationRuleClient.TopicsListKeys`

```go
ctx := context.TODO()
id := topicsauthorizationrule.NewTopicAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "topicName", "authorizationRuleName")

read, err := client.TopicsListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TopicsAuthorizationRuleClient.TopicsRegenerateKeys`

```go
ctx := context.TODO()
id := topicsauthorizationrule.NewTopicAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "topicName", "authorizationRuleName")

payload := topicsauthorizationrule.RegenerateAccessKeyParameters{
	// ...
}


read, err := client.TopicsRegenerateKeys(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
