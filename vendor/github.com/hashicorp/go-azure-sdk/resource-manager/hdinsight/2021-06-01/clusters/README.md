
## `github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/clusters` Documentation

The `clusters` SDK allows for interaction with the Azure Resource Manager Service `hdinsight` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/clusters"
```


### Client Initialization

```go
client := clusters.NewClustersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ClustersClient.Create`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := clusters.ClusterCreateParametersExtended{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.ExecuteScriptActions`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := clusters.ExecuteScriptActionParameters{
	// ...
}


if err := client.ExecuteScriptActionsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.Get`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ClustersClient.GetGatewaySettings`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

read, err := client.GetGatewaySettings(ctx, id)
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

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ClustersClient.ListByResourceGroup`

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


### Example Usage: `ClustersClient.Resize`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := clusters.ClusterResizeParameters{
	// ...
}


if err := client.ResizeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.RotateDiskEncryptionKey`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := clusters.ClusterDiskEncryptionParameters{
	// ...
}


if err := client.RotateDiskEncryptionKeyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.Update`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := clusters.ClusterPatchParameters{
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


### Example Usage: `ClustersClient.UpdateAutoScaleConfiguration`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := clusters.AutoscaleConfigurationUpdateParameter{
	// ...
}


if err := client.UpdateAutoScaleConfigurationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.UpdateGatewaySettings`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := clusters.UpdateGatewaySettingsParameters{
	// ...
}


if err := client.UpdateGatewaySettingsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ClustersClient.UpdateIdentityCertificate`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := clusters.UpdateClusterIdentityCertificateParameters{
	// ...
}


if err := client.UpdateIdentityCertificateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
