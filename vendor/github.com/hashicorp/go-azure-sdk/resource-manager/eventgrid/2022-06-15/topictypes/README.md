
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/topictypes` Documentation

The `topictypes` SDK allows for interaction with Azure Resource Manager `eventgrid` (API Version `2022-06-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/topictypes"
```


### Client Initialization

```go
client := topictypes.NewTopicTypesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TopicTypesClient.Get`

```go
ctx := context.TODO()
id := topictypes.NewTopicTypeID("topicTypeName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TopicTypesClient.List`

```go
ctx := context.TODO()


read, err := client.List(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TopicTypesClient.ListEventTypes`

```go
ctx := context.TODO()
id := topictypes.NewTopicTypeID("topicTypeName")

read, err := client.ListEventTypes(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
