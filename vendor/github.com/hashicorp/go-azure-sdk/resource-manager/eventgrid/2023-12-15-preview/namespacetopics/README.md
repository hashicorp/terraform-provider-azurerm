
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2023-12-15-preview/namespacetopics` Documentation

The `namespacetopics` SDK allows for interaction with Azure Resource Manager `eventgrid` (API Version `2023-12-15-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2023-12-15-preview/namespacetopics"
```


### Client Initialization

```go
client := namespacetopics.NewNamespaceTopicsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NamespaceTopicsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := namespacetopics.NewNamespaceTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "topicName")

payload := namespacetopics.NamespaceTopic{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NamespaceTopicsClient.Delete`

```go
ctx := context.TODO()
id := namespacetopics.NewNamespaceTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "topicName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NamespaceTopicsClient.Get`

```go
ctx := context.TODO()
id := namespacetopics.NewNamespaceTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "topicName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamespaceTopicsClient.ListByNamespace`

```go
ctx := context.TODO()
id := namespacetopics.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

// alternatively `client.ListByNamespace(ctx, id, namespacetopics.DefaultListByNamespaceOperationOptions())` can be used to do batched pagination
items, err := client.ListByNamespaceComplete(ctx, id, namespacetopics.DefaultListByNamespaceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NamespaceTopicsClient.ListSharedAccessKeys`

```go
ctx := context.TODO()
id := namespacetopics.NewNamespaceTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "topicName")

read, err := client.ListSharedAccessKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamespaceTopicsClient.RegenerateKey`

```go
ctx := context.TODO()
id := namespacetopics.NewNamespaceTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "topicName")

payload := namespacetopics.TopicRegenerateKeyRequest{
	// ...
}


if err := client.RegenerateKeyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NamespaceTopicsClient.Update`

```go
ctx := context.TODO()
id := namespacetopics.NewNamespaceTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "topicName")

payload := namespacetopics.NamespaceTopicUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
