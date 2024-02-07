
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/topics` Documentation

The `topics` SDK allows for interaction with the Azure Resource Manager Service `eventgrid` (API Version `2022-06-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/topics"
```


### Client Initialization

```go
client := topics.NewTopicsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TopicsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := topics.NewTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "topicValue")

payload := topics.Topic{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TopicsClient.Delete`

```go
ctx := context.TODO()
id := topics.NewTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "topicValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `TopicsClient.ExtensionTopicsGet`

```go
ctx := context.TODO()
id := topics.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.ExtensionTopicsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TopicsClient.Get`

```go
ctx := context.TODO()
id := topics.NewTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "topicValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TopicsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := topics.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, topics.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, topics.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TopicsClient.ListBySubscription`

```go
ctx := context.TODO()
id := topics.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, topics.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, topics.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TopicsClient.ListEventTypes`

```go
ctx := context.TODO()
id := topics.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.ListEventTypes(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TopicsClient.ListSharedAccessKeys`

```go
ctx := context.TODO()
id := topics.NewTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "topicValue")

read, err := client.ListSharedAccessKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TopicsClient.RegenerateKey`

```go
ctx := context.TODO()
id := topics.NewTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "topicValue")

payload := topics.TopicRegenerateKeyRequest{
	// ...
}


if err := client.RegenerateKeyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TopicsClient.Update`

```go
ctx := context.TODO()
id := topics.NewTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "topicValue")

payload := topics.TopicUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
