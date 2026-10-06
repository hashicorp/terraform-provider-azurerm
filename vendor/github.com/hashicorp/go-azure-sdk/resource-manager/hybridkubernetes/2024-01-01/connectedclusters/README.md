
## `github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters` Documentation

The `connectedclusters` SDK allows for interaction with Azure Resource Manager `hybridkubernetes` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
```


### Client Initialization

```go
client := connectedclusters.NewConnectedClustersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConnectedClustersClient.ConnectedClusterCreate`

```go
ctx := context.TODO()
id := connectedclusters.NewConnectedClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectedClusterName")

payload := connectedclusters.ConnectedCluster{
	// ...
}


if err := client.ConnectedClusterCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ConnectedClustersClient.ConnectedClusterDelete`

```go
ctx := context.TODO()
id := connectedclusters.NewConnectedClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectedClusterName")

if err := client.ConnectedClusterDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ConnectedClustersClient.ConnectedClusterGet`

```go
ctx := context.TODO()
id := connectedclusters.NewConnectedClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectedClusterName")

read, err := client.ConnectedClusterGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectedClustersClient.ConnectedClusterListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ConnectedClusterListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ConnectedClusterListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ConnectedClustersClient.ConnectedClusterListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ConnectedClusterListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ConnectedClusterListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ConnectedClustersClient.ConnectedClusterListClusterUserCredential`

```go
ctx := context.TODO()
id := connectedclusters.NewConnectedClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectedClusterName")

payload := connectedclusters.ListClusterUserCredentialProperties{
	// ...
}


read, err := client.ConnectedClusterListClusterUserCredential(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectedClustersClient.ConnectedClusterUpdate`

```go
ctx := context.TODO()
id := connectedclusters.NewConnectedClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectedClusterName")

payload := connectedclusters.ConnectedClusterPatch{
	// ...
}


read, err := client.ConnectedClusterUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
