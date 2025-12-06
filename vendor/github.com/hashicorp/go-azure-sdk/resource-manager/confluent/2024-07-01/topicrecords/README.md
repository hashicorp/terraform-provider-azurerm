
## `github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/topicrecords` Documentation

The `topicrecords` SDK allows for interaction with Azure Resource Manager `confluent` (API Version `2024-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/topicrecords"
```


### Client Initialization

```go
client := topicrecords.NewTopicRecordsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TopicRecordsClient.TopicsCreate`

```go
ctx := context.TODO()
id := topicrecords.NewTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId", "clusterId", "topicName")

payload := topicrecords.TopicRecord{
	// ...
}


read, err := client.TopicsCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TopicRecordsClient.TopicsDelete`

```go
ctx := context.TODO()
id := topicrecords.NewTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId", "clusterId", "topicName")

if err := client.TopicsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `TopicRecordsClient.TopicsGet`

```go
ctx := context.TODO()
id := topicrecords.NewTopicID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId", "clusterId", "topicName")

read, err := client.TopicsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TopicRecordsClient.TopicsList`

```go
ctx := context.TODO()
id := topicrecords.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId", "clusterId")

// alternatively `client.TopicsList(ctx, id, topicrecords.DefaultTopicsListOperationOptions())` can be used to do batched pagination
items, err := client.TopicsListComplete(ctx, id, topicrecords.DefaultTopicsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
