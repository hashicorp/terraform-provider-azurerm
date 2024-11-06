
## `github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/clusters` Documentation

The `clusters` SDK allows for interaction with Azure Resource Manager `kusto` (API Version `2023-08-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/clusters"
```


### Client Initialization

```go
client := clusters.NewClustersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ClustersClient.AddLanguageExtensions`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := clusters.LanguageExtensionsList{
	// ...
}


if err := client.AddLanguageExtensionsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := clusters.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := clusters.ClusterCheckNameRequest{
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


### Example Usage: `ClustersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := clusters.Cluster{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, clusters.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.DetachFollowerDatabases`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := clusters.FollowerDatabaseDefinition{
	// ...
}


if err := client.DetachFollowerDatabasesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.DiagnoseVirtualNetwork`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

if err := client.DiagnoseVirtualNetworkThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.Get`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

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
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ClustersClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.ListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ClustersClient.ListFollowerDatabases`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

read, err := client.ListFollowerDatabases(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ClustersClient.ListLanguageExtensions`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

read, err := client.ListLanguageExtensions(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ClustersClient.ListSkusByResource`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

read, err := client.ListSkusByResource(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ClustersClient.Migrate`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := clusters.ClusterMigrateRequest{
	// ...
}


if err := client.MigrateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.RemoveLanguageExtensions`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := clusters.LanguageExtensionsList{
	// ...
}


if err := client.RemoveLanguageExtensionsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.Start`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.Stop`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

if err := client.StopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.Update`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := clusters.ClusterUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload, clusters.DefaultUpdateOperationOptions()); err != nil {
	// handle the error
}
```
