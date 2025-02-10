
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/rules` Documentation

The `rules` SDK allows for interaction with Azure Resource Manager `servicebus` (API Version `2021-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/rules"
```


### Client Initialization

```go
client := rules.NewRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := rules.NewRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "topicName", "subscriptionName", "ruleName")

payload := rules.Rule{
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


### Example Usage: `RulesClient.Delete`

```go
ctx := context.TODO()
id := rules.NewRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "topicName", "subscriptionName", "ruleName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RulesClient.ListBySubscriptions`

```go
ctx := context.TODO()
id := rules.NewSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "topicName", "subscriptionName")

// alternatively `client.ListBySubscriptions(ctx, id, rules.DefaultListBySubscriptionsOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionsComplete(ctx, id, rules.DefaultListBySubscriptionsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
