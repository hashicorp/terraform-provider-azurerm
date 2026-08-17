
## `github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2025-09-01/mongoclusters` Documentation

The `mongoclusters` SDK allows for interaction with Azure Resource Manager `mongocluster` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2025-09-01/mongoclusters"
```


### Client Initialization

```go
client := mongoclusters.NewMongoClustersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MongoClustersClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := mongoclusters.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := mongoclusters.CheckNameAvailabilityRequest{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MongoClustersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := mongoclusters.NewMongoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mongoClusterName")

payload := mongoclusters.MongoCluster{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `MongoClustersClient.Delete`

```go
ctx := context.TODO()
id := mongoclusters.NewMongoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mongoClusterName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `MongoClustersClient.Get`

```go
ctx := context.TODO()
id := mongoclusters.NewMongoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mongoClusterName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MongoClustersClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MongoClustersClient.ListByResourceGroup`

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


### Example Usage: `MongoClustersClient.ListConnectionStrings`

```go
ctx := context.TODO()
id := mongoclusters.NewMongoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mongoClusterName")

read, err := client.ListConnectionStrings(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MongoClustersClient.Promote`

```go
ctx := context.TODO()
id := mongoclusters.NewMongoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mongoClusterName")

payload := mongoclusters.PromoteReplicaRequest{
	// ...
}


if err := client.PromoteThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `MongoClustersClient.Update`

```go
ctx := context.TODO()
id := mongoclusters.NewMongoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mongoClusterName")

payload := mongoclusters.MongoClusterUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
