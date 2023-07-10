
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/publishers` Documentation

The `publishers` SDK allows for interaction with the Azure Resource Manager Service `azurestackhci` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/publishers"
```


### Client Initialization

```go
client := publishers.NewPublishersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PublishersClient.PublishersGet`

```go
ctx := context.TODO()
id := publishers.NewPublisherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "publisherValue")

read, err := client.PublishersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PublishersClient.PublishersListByCluster`

```go
ctx := context.TODO()
id := publishers.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

// alternatively `client.PublishersListByCluster(ctx, id)` can be used to do batched pagination
items, err := client.PublishersListByClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
