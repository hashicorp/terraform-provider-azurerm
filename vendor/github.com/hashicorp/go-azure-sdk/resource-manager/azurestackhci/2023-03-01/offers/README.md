
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/offers` Documentation

The `offers` SDK allows for interaction with the Azure Resource Manager Service `azurestackhci` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/offers"
```


### Client Initialization

```go
client := offers.NewOffersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OffersClient.Get`

```go
ctx := context.TODO()
id := offers.NewOfferID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "publisherValue", "offerValue")

read, err := client.Get(ctx, id, offers.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OffersClient.ListByCluster`

```go
ctx := context.TODO()
id := offers.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

// alternatively `client.ListByCluster(ctx, id, offers.DefaultListByClusterOperationOptions())` can be used to do batched pagination
items, err := client.ListByClusterComplete(ctx, id, offers.DefaultListByClusterOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OffersClient.ListByPublisher`

```go
ctx := context.TODO()
id := offers.NewPublisherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "publisherValue")

// alternatively `client.ListByPublisher(ctx, id, offers.DefaultListByPublisherOperationOptions())` can be used to do batched pagination
items, err := client.ListByPublisherComplete(ctx, id, offers.DefaultListByPublisherOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
