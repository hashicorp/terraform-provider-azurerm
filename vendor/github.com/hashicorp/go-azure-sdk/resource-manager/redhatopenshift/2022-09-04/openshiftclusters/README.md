
## `github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2022-09-04/openshiftclusters` Documentation

The `openshiftclusters` SDK allows for interaction with the Azure Resource Manager Service `redhatopenshift` (API Version `2022-09-04`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2022-09-04/openshiftclusters"
```


### Client Initialization

```go
client := openshiftclusters.NewOpenShiftClustersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OpenShiftClustersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := openshiftclusters.NewProviderOpenShiftClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "openShiftClusterValue")

payload := openshiftclusters.OpenShiftCluster{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenShiftClustersClient.Delete`

```go
ctx := context.TODO()
id := openshiftclusters.NewProviderOpenShiftClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "openShiftClusterValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenShiftClustersClient.Get`

```go
ctx := context.TODO()
id := openshiftclusters.NewProviderOpenShiftClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "openShiftClusterValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenShiftClustersClient.List`

```go
ctx := context.TODO()
id := openshiftclusters.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenShiftClustersClient.ListAdminCredentials`

```go
ctx := context.TODO()
id := openshiftclusters.NewProviderOpenShiftClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "openShiftClusterValue")

read, err := client.ListAdminCredentials(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenShiftClustersClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := openshiftclusters.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenShiftClustersClient.ListCredentials`

```go
ctx := context.TODO()
id := openshiftclusters.NewProviderOpenShiftClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "openShiftClusterValue")

read, err := client.ListCredentials(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenShiftClustersClient.Update`

```go
ctx := context.TODO()
id := openshiftclusters.NewProviderOpenShiftClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "openShiftClusterValue")

payload := openshiftclusters.OpenShiftClusterUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
