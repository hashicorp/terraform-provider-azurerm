
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubsclusters` Documentation

The `eventhubsclusters` SDK allows for interaction with Azure Resource Manager `eventhub` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubsclusters"
```


### Client Initialization

```go
client := eventhubsclusters.NewEventHubsClustersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EventHubsClustersClient.ClustersCreateOrUpdate`

```go
ctx := context.TODO()
id := eventhubsclusters.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := eventhubsclusters.Cluster{
	// ...
}


if err := client.ClustersCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EventHubsClustersClient.ClustersDelete`

```go
ctx := context.TODO()
id := eventhubsclusters.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

if err := client.ClustersDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EventHubsClustersClient.ClustersGet`

```go
ctx := context.TODO()
id := eventhubsclusters.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

read, err := client.ClustersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventHubsClustersClient.ClustersListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ClustersListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ClustersListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventHubsClustersClient.ClustersListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ClustersListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ClustersListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EventHubsClustersClient.ClustersUpdate`

```go
ctx := context.TODO()
id := eventhubsclusters.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := eventhubsclusters.Cluster{
	// ...
}


if err := client.ClustersUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
