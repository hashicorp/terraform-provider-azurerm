
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/managedclusters` Documentation

The `managedclusters` SDK allows for interaction with the Azure Resource Manager Service `containerservice` (API Version `2023-03-02-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/managedclusters"
```


### Client Initialization

```go
client := managedclusters.NewManagedClustersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedClustersClient.AbortLatestOperation`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

if err := client.AbortLatestOperationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedClustersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

payload := managedclusters.ManagedCluster{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedClustersClient.Delete`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

if err := client.DeleteThenPoll(ctx, id, managedclusters.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedClustersClient.Get`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedClustersClient.GetAccessProfile`

```go
ctx := context.TODO()
id := managedclusters.NewAccessProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue", "accessProfileValue")

read, err := client.GetAccessProfile(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedClustersClient.GetCommandResult`

```go
ctx := context.TODO()
id := managedclusters.NewCommandResultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue", "commandIdValue")

read, err := client.GetCommandResult(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedClustersClient.GetOSOptions`

```go
ctx := context.TODO()
id := managedclusters.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

read, err := client.GetOSOptions(ctx, id, managedclusters.DefaultGetOSOptionsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedClustersClient.GetUpgradeProfile`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

read, err := client.GetUpgradeProfile(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedClustersClient.List`

```go
ctx := context.TODO()
id := managedclusters.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedClustersClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := managedclusters.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedClustersClient.ListClusterAdminCredentials`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

read, err := client.ListClusterAdminCredentials(ctx, id, managedclusters.DefaultListClusterAdminCredentialsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedClustersClient.ListClusterMonitoringUserCredentials`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

read, err := client.ListClusterMonitoringUserCredentials(ctx, id, managedclusters.DefaultListClusterMonitoringUserCredentialsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedClustersClient.ListClusterUserCredentials`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

read, err := client.ListClusterUserCredentials(ctx, id, managedclusters.DefaultListClusterUserCredentialsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedClustersClient.ListKubernetesVersions`

```go
ctx := context.TODO()
id := managedclusters.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

read, err := client.ListKubernetesVersions(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedClustersClient.ListOutboundNetworkDependenciesEndpoints`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

// alternatively `client.ListOutboundNetworkDependenciesEndpoints(ctx, id)` can be used to do batched pagination
items, err := client.ListOutboundNetworkDependenciesEndpointsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedClustersClient.ResetAADProfile`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

payload := managedclusters.ManagedClusterAADProfile{
	// ...
}


if err := client.ResetAADProfileThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedClustersClient.ResetServicePrincipalProfile`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

payload := managedclusters.ManagedClusterServicePrincipalProfile{
	// ...
}


if err := client.ResetServicePrincipalProfileThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedClustersClient.RotateClusterCertificates`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

if err := client.RotateClusterCertificatesThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedClustersClient.RotateServiceAccountSigningKeys`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

if err := client.RotateServiceAccountSigningKeysThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedClustersClient.RunCommand`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

payload := managedclusters.RunCommandRequest{
	// ...
}


if err := client.RunCommandThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedClustersClient.Start`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedClustersClient.Stop`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

if err := client.StopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedClustersClient.UpdateTags`

```go
ctx := context.TODO()
id := managedclusters.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

payload := managedclusters.TagsObject{
	// ...
}


if err := client.UpdateTagsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
