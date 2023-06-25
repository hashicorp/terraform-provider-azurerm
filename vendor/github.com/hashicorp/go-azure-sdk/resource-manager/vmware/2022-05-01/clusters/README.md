
## `github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/clusters` Documentation

The `clusters` SDK allows for interaction with the Azure Resource Manager Service `vmware` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/clusters"
```


### Client Initialization

```go
client := clusters.NewClustersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ClustersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := clusters.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudValue", "clusterValue")

payload := clusters.Cluster{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.Delete`

```go
ctx := context.TODO()
id := clusters.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudValue", "clusterValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.Get`

```go
ctx := context.TODO()
id := clusters.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudValue", "clusterValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ClustersClient.List`

```go
ctx := context.TODO()
id := clusters.NewPrivateCloudID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ClustersClient.Update`

```go
ctx := context.TODO()
id := clusters.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudValue", "clusterValue")

payload := clusters.ClusterUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
