
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/domaintopics` Documentation

The `domaintopics` SDK allows for interaction with the Azure Resource Manager Service `eventgrid` (API Version `2022-06-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/domaintopics"
```


### Client Initialization

```go
client := domaintopics.NewDomainTopicsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DomainTopicsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := domaintopics.NewDomainTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "topicValue")

if err := client.CreateOrUpdateThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DomainTopicsClient.Delete`

```go
ctx := context.TODO()
id := domaintopics.NewDomainTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "topicValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DomainTopicsClient.Get`

```go
ctx := context.TODO()
id := domaintopics.NewDomainTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue", "topicValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DomainTopicsClient.ListByDomain`

```go
ctx := context.TODO()
id := domaintopics.NewDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "domainValue")

// alternatively `client.ListByDomain(ctx, id, domaintopics.DefaultListByDomainOperationOptions())` can be used to do batched pagination
items, err := client.ListByDomainComplete(ctx, id, domaintopics.DefaultListByDomainOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
