
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/subscriptions` Documentation

The `subscriptions` SDK allows for interaction with the Azure Resource Manager Service `servicebus` (API Version `2021-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/subscriptions"
```


### Client Initialization

```go
client := subscriptions.NewSubscriptionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SubscriptionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := subscriptions.NewSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "topicValue", "subscriptionValue")

payload := subscriptions.SBSubscription{
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


### Example Usage: `SubscriptionsClient.Delete`

```go
ctx := context.TODO()
id := subscriptions.NewSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "topicValue", "subscriptionValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionsClient.Get`

```go
ctx := context.TODO()
id := subscriptions.NewSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "topicValue", "subscriptionValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionsClient.ListByTopic`

```go
ctx := context.TODO()
id := subscriptions.NewTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "topicValue")

// alternatively `client.ListByTopic(ctx, id, subscriptions.DefaultListByTopicOperationOptions())` can be used to do batched pagination
items, err := client.ListByTopicComplete(ctx, id, subscriptions.DefaultListByTopicOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SubscriptionsClient.RulesGet`

```go
ctx := context.TODO()
id := subscriptions.NewRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "topicValue", "subscriptionValue", "ruleValue")

read, err := client.RulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
