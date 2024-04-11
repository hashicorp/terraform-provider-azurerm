
## `github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionendpoints` Documentation

The `datacollectionendpoints` SDK allows for interaction with the Azure Resource Manager Service `insights` (API Version `2022-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionendpoints"
```


### Client Initialization

```go
client := datacollectionendpoints.NewDataCollectionEndpointsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataCollectionEndpointsClient.Create`

```go
ctx := context.TODO()
id := datacollectionendpoints.NewDataCollectionEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataCollectionEndpointValue")

payload := datacollectionendpoints.DataCollectionEndpointResource{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataCollectionEndpointsClient.Delete`

```go
ctx := context.TODO()
id := datacollectionendpoints.NewDataCollectionEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataCollectionEndpointValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataCollectionEndpointsClient.Get`

```go
ctx := context.TODO()
id := datacollectionendpoints.NewDataCollectionEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataCollectionEndpointValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataCollectionEndpointsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DataCollectionEndpointsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DataCollectionEndpointsClient.Update`

```go
ctx := context.TODO()
id := datacollectionendpoints.NewDataCollectionEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataCollectionEndpointValue")

payload := datacollectionendpoints.ResourceForUpdate{
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
