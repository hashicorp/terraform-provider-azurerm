
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudvmclusters` Documentation

The `cloudvmclusters` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2024-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudvmclusters"
```


### Client Initialization

```go
client := cloudvmclusters.NewCloudVMClustersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CloudVMClustersClient.AddVMs`

```go
ctx := context.TODO()
id := cloudvmclusters.NewCloudVMClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudVmClusterName")

payload := cloudvmclusters.AddRemoveDbNode{
	// ...
}


if err := client.AddVMsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudVMClustersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := cloudvmclusters.NewCloudVMClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudVmClusterName")

payload := cloudvmclusters.CloudVMCluster{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudVMClustersClient.Delete`

```go
ctx := context.TODO()
id := cloudvmclusters.NewCloudVMClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudVmClusterName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CloudVMClustersClient.Get`

```go
ctx := context.TODO()
id := cloudvmclusters.NewCloudVMClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudVmClusterName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CloudVMClustersClient.ListByResourceGroup`

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


### Example Usage: `CloudVMClustersClient.ListBySubscription`

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


### Example Usage: `CloudVMClustersClient.ListPrivateIPAddresses`

```go
ctx := context.TODO()
id := cloudvmclusters.NewCloudVMClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudVmClusterName")

payload := cloudvmclusters.PrivateIPAddressesFilter{
	// ...
}


read, err := client.ListPrivateIPAddresses(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CloudVMClustersClient.RemoveVMs`

```go
ctx := context.TODO()
id := cloudvmclusters.NewCloudVMClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudVmClusterName")

payload := cloudvmclusters.AddRemoveDbNode{
	// ...
}


if err := client.RemoveVMsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudVMClustersClient.Update`

```go
ctx := context.TODO()
id := cloudvmclusters.NewCloudVMClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudVmClusterName")

payload := cloudvmclusters.CloudVMClusterUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
