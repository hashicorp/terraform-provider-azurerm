
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/exadbvmclusters` Documentation

The `exadbvmclusters` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/exadbvmclusters"
```


### Client Initialization

```go
client := exadbvmclusters.NewExadbVMClustersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExadbVMClustersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := exadbvmclusters.NewExadbVMClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "exadbVmClusterName")

payload := exadbvmclusters.ExadbVMCluster{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExadbVMClustersClient.Delete`

```go
ctx := context.TODO()
id := exadbvmclusters.NewExadbVMClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "exadbVmClusterName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExadbVMClustersClient.Get`

```go
ctx := context.TODO()
id := exadbvmclusters.NewExadbVMClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "exadbVmClusterName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExadbVMClustersClient.ListByResourceGroup`

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


### Example Usage: `ExadbVMClustersClient.ListBySubscription`

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


### Example Usage: `ExadbVMClustersClient.RemoveVMs`

```go
ctx := context.TODO()
id := exadbvmclusters.NewExadbVMClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "exadbVmClusterName")

payload := exadbvmclusters.RemoveVirtualMachineFromExadbVMClusterDetails{
	// ...
}


if err := client.RemoveVMsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExadbVMClustersClient.Update`

```go
ctx := context.TODO()
id := exadbvmclusters.NewExadbVMClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "exadbVmClusterName")

payload := exadbvmclusters.ExadbVMClusterUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
