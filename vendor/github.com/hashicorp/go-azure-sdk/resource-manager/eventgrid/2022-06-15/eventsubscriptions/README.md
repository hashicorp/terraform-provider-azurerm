
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/eventsubscriptions` Documentation

The `eventsubscriptions` SDK allows for interaction with the Azure Resource Manager Service `eventgrid` (API Version `2022-06-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/eventsubscriptions"
```


### Client Initialization

```go
client := eventsubscriptions.NewEventSubscriptionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EventSubscriptionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := eventsubscriptions.NewScopedEventSubscriptionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "eventSubscriptionValue")

payload := eventsubscriptions.EventSubscription{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.Delete`

```go
ctx := context.TODO()
id := eventsubscriptions.NewScopedEventSubscriptionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "eventSubscriptionValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.DomainEventSubscriptionsCreateOrUpdate`

```go
ctx := context.TODO()
id := eventsubscriptions.NewDomainEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "eventSubscriptionValue")

payload := eventsubscriptions.EventSubscription{
	// ...
}


if err := client.DomainEventSubscriptionsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.DomainEventSubscriptionsDelete`

```go
ctx := context.TODO()
id := eventsubscriptions.NewDomainEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "eventSubscriptionValue")

if err := client.DomainEventSubscriptionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.DomainEventSubscriptionsGet`

```go
ctx := context.TODO()
id := eventsubscriptions.NewDomainEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "eventSubscriptionValue")

read, err := client.DomainEventSubscriptionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.DomainEventSubscriptionsGetDeliveryAttributes`

```go
ctx := context.TODO()
id := eventsubscriptions.NewDomainEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "eventSubscriptionValue")

read, err := client.DomainEventSubscriptionsGetDeliveryAttributes(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.DomainEventSubscriptionsGetFullUrl`

```go
ctx := context.TODO()
id := eventsubscriptions.NewDomainEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "eventSubscriptionValue")

read, err := client.DomainEventSubscriptionsGetFullUrl(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.DomainEventSubscriptionsList`

```go
ctx := context.TODO()
id := eventsubscriptions.NewDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue")

// alternatively `client.DomainEventSubscriptionsList(ctx, id, eventsubscriptions.DefaultDomainEventSubscriptionsListOperationOptions())` can be used to do batched pagination
items, err := client.DomainEventSubscriptionsListComplete(ctx, id, eventsubscriptions.DefaultDomainEventSubscriptionsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventSubscriptionsClient.DomainEventSubscriptionsUpdate`

```go
ctx := context.TODO()
id := eventsubscriptions.NewDomainEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "eventSubscriptionValue")

payload := eventsubscriptions.EventSubscriptionUpdateParameters{
	// ...
}


if err := client.DomainEventSubscriptionsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.DomainTopicEventSubscriptionsCreateOrUpdate`

```go
ctx := context.TODO()
id := eventsubscriptions.NewTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "topicValue", "eventSubscriptionValue")

payload := eventsubscriptions.EventSubscription{
	// ...
}


if err := client.DomainTopicEventSubscriptionsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.DomainTopicEventSubscriptionsDelete`

```go
ctx := context.TODO()
id := eventsubscriptions.NewTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "topicValue", "eventSubscriptionValue")

if err := client.DomainTopicEventSubscriptionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.DomainTopicEventSubscriptionsGet`

```go
ctx := context.TODO()
id := eventsubscriptions.NewTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "topicValue", "eventSubscriptionValue")

read, err := client.DomainTopicEventSubscriptionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.DomainTopicEventSubscriptionsGetDeliveryAttributes`

```go
ctx := context.TODO()
id := eventsubscriptions.NewTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "topicValue", "eventSubscriptionValue")

read, err := client.DomainTopicEventSubscriptionsGetDeliveryAttributes(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.DomainTopicEventSubscriptionsGetFullUrl`

```go
ctx := context.TODO()
id := eventsubscriptions.NewTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "topicValue", "eventSubscriptionValue")

read, err := client.DomainTopicEventSubscriptionsGetFullUrl(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.DomainTopicEventSubscriptionsList`

```go
ctx := context.TODO()
id := eventsubscriptions.NewDomainTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "topicValue")

// alternatively `client.DomainTopicEventSubscriptionsList(ctx, id, eventsubscriptions.DefaultDomainTopicEventSubscriptionsListOperationOptions())` can be used to do batched pagination
items, err := client.DomainTopicEventSubscriptionsListComplete(ctx, id, eventsubscriptions.DefaultDomainTopicEventSubscriptionsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventSubscriptionsClient.DomainTopicEventSubscriptionsUpdate`

```go
ctx := context.TODO()
id := eventsubscriptions.NewTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "topicValue", "eventSubscriptionValue")

payload := eventsubscriptions.EventSubscriptionUpdateParameters{
	// ...
}


if err := client.DomainTopicEventSubscriptionsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.Get`

```go
ctx := context.TODO()
id := eventsubscriptions.NewScopedEventSubscriptionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "eventSubscriptionValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.GetDeliveryAttributes`

```go
ctx := context.TODO()
id := eventsubscriptions.NewScopedEventSubscriptionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "eventSubscriptionValue")

read, err := client.GetDeliveryAttributes(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.GetFullUrl`

```go
ctx := context.TODO()
id := eventsubscriptions.NewScopedEventSubscriptionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "eventSubscriptionValue")

read, err := client.GetFullUrl(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.ListByDomainTopic`

```go
ctx := context.TODO()
id := eventsubscriptions.NewDomainTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "topicValue")

// alternatively `client.ListByDomainTopic(ctx, id, eventsubscriptions.DefaultListByDomainTopicOperationOptions())` can be used to do batched pagination
items, err := client.ListByDomainTopicComplete(ctx, id, eventsubscriptions.DefaultListByDomainTopicOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventSubscriptionsClient.ListByResource`

```go
ctx := context.TODO()
id := eventsubscriptions.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ListByResource(ctx, id, eventsubscriptions.DefaultListByResourceOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceComplete(ctx, id, eventsubscriptions.DefaultListByResourceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventSubscriptionsClient.ListGlobalByResourceGroup`

```go
ctx := context.TODO()
id := eventsubscriptions.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListGlobalByResourceGroup(ctx, id, eventsubscriptions.DefaultListGlobalByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListGlobalByResourceGroupComplete(ctx, id, eventsubscriptions.DefaultListGlobalByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventSubscriptionsClient.ListGlobalByResourceGroupForTopicType`

```go
ctx := context.TODO()
id := eventsubscriptions.NewResourceGroupProviderTopicTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "topicTypeValue")

// alternatively `client.ListGlobalByResourceGroupForTopicType(ctx, id, eventsubscriptions.DefaultListGlobalByResourceGroupForTopicTypeOperationOptions())` can be used to do batched pagination
items, err := client.ListGlobalByResourceGroupForTopicTypeComplete(ctx, id, eventsubscriptions.DefaultListGlobalByResourceGroupForTopicTypeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventSubscriptionsClient.ListGlobalBySubscription`

```go
ctx := context.TODO()
id := eventsubscriptions.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListGlobalBySubscription(ctx, id, eventsubscriptions.DefaultListGlobalBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListGlobalBySubscriptionComplete(ctx, id, eventsubscriptions.DefaultListGlobalBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventSubscriptionsClient.ListGlobalBySubscriptionForTopicType`

```go
ctx := context.TODO()
id := eventsubscriptions.NewProviderTopicTypeID("12345678-1234-9876-4563-123456789012", "topicTypeValue")

// alternatively `client.ListGlobalBySubscriptionForTopicType(ctx, id, eventsubscriptions.DefaultListGlobalBySubscriptionForTopicTypeOperationOptions())` can be used to do batched pagination
items, err := client.ListGlobalBySubscriptionForTopicTypeComplete(ctx, id, eventsubscriptions.DefaultListGlobalBySubscriptionForTopicTypeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventSubscriptionsClient.ListRegionalByResourceGroup`

```go
ctx := context.TODO()
id := eventsubscriptions.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationValue")

// alternatively `client.ListRegionalByResourceGroup(ctx, id, eventsubscriptions.DefaultListRegionalByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListRegionalByResourceGroupComplete(ctx, id, eventsubscriptions.DefaultListRegionalByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventSubscriptionsClient.ListRegionalByResourceGroupForTopicType`

```go
ctx := context.TODO()
id := eventsubscriptions.NewProviderLocationTopicTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationValue", "topicTypeValue")

// alternatively `client.ListRegionalByResourceGroupForTopicType(ctx, id, eventsubscriptions.DefaultListRegionalByResourceGroupForTopicTypeOperationOptions())` can be used to do batched pagination
items, err := client.ListRegionalByResourceGroupForTopicTypeComplete(ctx, id, eventsubscriptions.DefaultListRegionalByResourceGroupForTopicTypeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventSubscriptionsClient.ListRegionalBySubscription`

```go
ctx := context.TODO()
id := eventsubscriptions.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

// alternatively `client.ListRegionalBySubscription(ctx, id, eventsubscriptions.DefaultListRegionalBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListRegionalBySubscriptionComplete(ctx, id, eventsubscriptions.DefaultListRegionalBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventSubscriptionsClient.ListRegionalBySubscriptionForTopicType`

```go
ctx := context.TODO()
id := eventsubscriptions.NewLocationTopicTypeID("12345678-1234-9876-4563-123456789012", "locationValue", "topicTypeValue")

// alternatively `client.ListRegionalBySubscriptionForTopicType(ctx, id, eventsubscriptions.DefaultListRegionalBySubscriptionForTopicTypeOperationOptions())` can be used to do batched pagination
items, err := client.ListRegionalBySubscriptionForTopicTypeComplete(ctx, id, eventsubscriptions.DefaultListRegionalBySubscriptionForTopicTypeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventSubscriptionsClient.PartnerTopicEventSubscriptionsCreateOrUpdate`

```go
ctx := context.TODO()
id := eventsubscriptions.NewPartnerTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerTopicValue", "eventSubscriptionValue")

payload := eventsubscriptions.EventSubscription{
	// ...
}


if err := client.PartnerTopicEventSubscriptionsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.PartnerTopicEventSubscriptionsDelete`

```go
ctx := context.TODO()
id := eventsubscriptions.NewPartnerTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerTopicValue", "eventSubscriptionValue")

if err := client.PartnerTopicEventSubscriptionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.PartnerTopicEventSubscriptionsGet`

```go
ctx := context.TODO()
id := eventsubscriptions.NewPartnerTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerTopicValue", "eventSubscriptionValue")

read, err := client.PartnerTopicEventSubscriptionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.PartnerTopicEventSubscriptionsGetDeliveryAttributes`

```go
ctx := context.TODO()
id := eventsubscriptions.NewPartnerTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerTopicValue", "eventSubscriptionValue")

read, err := client.PartnerTopicEventSubscriptionsGetDeliveryAttributes(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.PartnerTopicEventSubscriptionsGetFullUrl`

```go
ctx := context.TODO()
id := eventsubscriptions.NewPartnerTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerTopicValue", "eventSubscriptionValue")

read, err := client.PartnerTopicEventSubscriptionsGetFullUrl(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.PartnerTopicEventSubscriptionsListByPartnerTopic`

```go
ctx := context.TODO()
id := eventsubscriptions.NewPartnerTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerTopicValue")

// alternatively `client.PartnerTopicEventSubscriptionsListByPartnerTopic(ctx, id, eventsubscriptions.DefaultPartnerTopicEventSubscriptionsListByPartnerTopicOperationOptions())` can be used to do batched pagination
items, err := client.PartnerTopicEventSubscriptionsListByPartnerTopicComplete(ctx, id, eventsubscriptions.DefaultPartnerTopicEventSubscriptionsListByPartnerTopicOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventSubscriptionsClient.PartnerTopicEventSubscriptionsUpdate`

```go
ctx := context.TODO()
id := eventsubscriptions.NewPartnerTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerTopicValue", "eventSubscriptionValue")

payload := eventsubscriptions.EventSubscriptionUpdateParameters{
	// ...
}


if err := client.PartnerTopicEventSubscriptionsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.SystemTopicEventSubscriptionsCreateOrUpdate`

```go
ctx := context.TODO()
id := eventsubscriptions.NewSystemTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "systemTopicValue", "eventSubscriptionValue")

payload := eventsubscriptions.EventSubscription{
	// ...
}


if err := client.SystemTopicEventSubscriptionsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.SystemTopicEventSubscriptionsDelete`

```go
ctx := context.TODO()
id := eventsubscriptions.NewSystemTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "systemTopicValue", "eventSubscriptionValue")

if err := client.SystemTopicEventSubscriptionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.SystemTopicEventSubscriptionsGet`

```go
ctx := context.TODO()
id := eventsubscriptions.NewSystemTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "systemTopicValue", "eventSubscriptionValue")

read, err := client.SystemTopicEventSubscriptionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.SystemTopicEventSubscriptionsGetDeliveryAttributes`

```go
ctx := context.TODO()
id := eventsubscriptions.NewSystemTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "systemTopicValue", "eventSubscriptionValue")

read, err := client.SystemTopicEventSubscriptionsGetDeliveryAttributes(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.SystemTopicEventSubscriptionsGetFullUrl`

```go
ctx := context.TODO()
id := eventsubscriptions.NewSystemTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "systemTopicValue", "eventSubscriptionValue")

read, err := client.SystemTopicEventSubscriptionsGetFullUrl(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.SystemTopicEventSubscriptionsListBySystemTopic`

```go
ctx := context.TODO()
id := eventsubscriptions.NewSystemTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "systemTopicValue")

// alternatively `client.SystemTopicEventSubscriptionsListBySystemTopic(ctx, id, eventsubscriptions.DefaultSystemTopicEventSubscriptionsListBySystemTopicOperationOptions())` can be used to do batched pagination
items, err := client.SystemTopicEventSubscriptionsListBySystemTopicComplete(ctx, id, eventsubscriptions.DefaultSystemTopicEventSubscriptionsListBySystemTopicOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventSubscriptionsClient.SystemTopicEventSubscriptionsUpdate`

```go
ctx := context.TODO()
id := eventsubscriptions.NewSystemTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "systemTopicValue", "eventSubscriptionValue")

payload := eventsubscriptions.EventSubscriptionUpdateParameters{
	// ...
}


if err := client.SystemTopicEventSubscriptionsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.TopicEventSubscriptionsCreateOrUpdate`

```go
ctx := context.TODO()
id := eventsubscriptions.NewEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "topicValue", "eventSubscriptionValue")

payload := eventsubscriptions.EventSubscription{
	// ...
}


if err := client.TopicEventSubscriptionsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.TopicEventSubscriptionsDelete`

```go
ctx := context.TODO()
id := eventsubscriptions.NewEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "topicValue", "eventSubscriptionValue")

if err := client.TopicEventSubscriptionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.TopicEventSubscriptionsGet`

```go
ctx := context.TODO()
id := eventsubscriptions.NewEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "topicValue", "eventSubscriptionValue")

read, err := client.TopicEventSubscriptionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.TopicEventSubscriptionsGetDeliveryAttributes`

```go
ctx := context.TODO()
id := eventsubscriptions.NewEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "topicValue", "eventSubscriptionValue")

read, err := client.TopicEventSubscriptionsGetDeliveryAttributes(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.TopicEventSubscriptionsGetFullUrl`

```go
ctx := context.TODO()
id := eventsubscriptions.NewEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "topicValue", "eventSubscriptionValue")

read, err := client.TopicEventSubscriptionsGetFullUrl(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventSubscriptionsClient.TopicEventSubscriptionsList`

```go
ctx := context.TODO()
id := eventsubscriptions.NewTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "topicValue")

// alternatively `client.TopicEventSubscriptionsList(ctx, id, eventsubscriptions.DefaultTopicEventSubscriptionsListOperationOptions())` can be used to do batched pagination
items, err := client.TopicEventSubscriptionsListComplete(ctx, id, eventsubscriptions.DefaultTopicEventSubscriptionsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventSubscriptionsClient.TopicEventSubscriptionsUpdate`

```go
ctx := context.TODO()
id := eventsubscriptions.NewEventSubscriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "topicValue", "eventSubscriptionValue")

payload := eventsubscriptions.EventSubscriptionUpdateParameters{
	// ...
}


if err := client.TopicEventSubscriptionsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EventSubscriptionsClient.Update`

```go
ctx := context.TODO()
id := eventsubscriptions.NewScopedEventSubscriptionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "eventSubscriptionValue")

payload := eventsubscriptions.EventSubscriptionUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
