
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/systemtopics` Documentation

The `systemtopics` SDK allows for interaction with the Azure Resource Manager Service `eventgrid` (API Version `2022-06-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/systemtopics"
```


### Client Initialization

```go
client := systemtopics.NewSystemTopicsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SystemTopicsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := systemtopics.NewSystemTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "systemTopicValue")

payload := systemtopics.SystemTopic{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SystemTopicsClient.Delete`

```go
ctx := context.TODO()
id := systemtopics.NewSystemTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "systemTopicValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SystemTopicsClient.Get`

```go
ctx := context.TODO()
id := systemtopics.NewSystemTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "systemTopicValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SystemTopicsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := systemtopics.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, systemtopics.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, systemtopics.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SystemTopicsClient.ListBySubscription`

```go
ctx := context.TODO()
id := systemtopics.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, systemtopics.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, systemtopics.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SystemTopicsClient.Update`

```go
ctx := context.TODO()
id := systemtopics.NewSystemTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "systemTopicValue")

payload := systemtopics.SystemTopicUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
