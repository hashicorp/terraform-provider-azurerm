
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/publishers` Documentation

The `publishers` SDK allows for interaction with Azure Resource Manager `azurestackhci` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/publishers"
```


### Client Initialization

```go
client := publishers.NewPublishersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PublishersClient.Get`

```go
ctx := context.TODO()
id := publishers.NewPublisherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "publisherName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PublishersClient.ListByCluster`

```go
ctx := context.TODO()
id := publishers.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

// alternatively `client.ListByCluster(ctx, id)` can be used to do batched pagination
items, err := client.ListByClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
