
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnertopics` Documentation

The `partnertopics` SDK allows for interaction with Azure Resource Manager `eventgrid` (API Version `2022-06-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnertopics"
```


### Client Initialization

```go
client := partnertopics.NewPartnerTopicsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PartnerTopicsClient.Activate`

```go
ctx := context.TODO()
id := partnertopics.NewPartnerTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerTopicName")

read, err := client.Activate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PartnerTopicsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := partnertopics.NewPartnerTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerTopicName")

payload := partnertopics.PartnerTopic{
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


### Example Usage: `PartnerTopicsClient.Deactivate`

```go
ctx := context.TODO()
id := partnertopics.NewPartnerTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerTopicName")

read, err := client.Deactivate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PartnerTopicsClient.Delete`

```go
ctx := context.TODO()
id := partnertopics.NewPartnerTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerTopicName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PartnerTopicsClient.Get`

```go
ctx := context.TODO()
id := partnertopics.NewPartnerTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerTopicName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PartnerTopicsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, partnertopics.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, partnertopics.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PartnerTopicsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, partnertopics.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, partnertopics.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PartnerTopicsClient.Update`

```go
ctx := context.TODO()
id := partnertopics.NewPartnerTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerTopicName")

payload := partnertopics.PartnerTopicUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
