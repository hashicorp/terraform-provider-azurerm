
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


### Example Usage: `OffersClient.OffersGet`

```go
ctx := context.TODO()
id := offers.NewOfferID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "publisherValue", "offerValue")

read, err := client.OffersGet(ctx, id, offers.DefaultOffersGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OffersClient.OffersListByCluster`

```go
ctx := context.TODO()
id := offers.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

// alternatively `client.OffersListByCluster(ctx, id, offers.DefaultOffersListByClusterOperationOptions())` can be used to do batched pagination
items, err := client.OffersListByClusterComplete(ctx, id, offers.DefaultOffersListByClusterOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OffersClient.OffersListByPublisher`

```go
ctx := context.TODO()
id := offers.NewPublisherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "publisherValue")

// alternatively `client.OffersListByPublisher(ctx, id, offers.DefaultOffersListByPublisherOperationOptions())` can be used to do batched pagination
items, err := client.OffersListByPublisherComplete(ctx, id, offers.DefaultOffersListByPublisherOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
