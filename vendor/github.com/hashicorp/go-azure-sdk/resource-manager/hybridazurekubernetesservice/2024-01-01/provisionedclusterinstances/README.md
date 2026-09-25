
## `github.com/hashicorp/go-azure-sdk/resource-manager/hybridazurekubernetesservice/2024-01-01/provisionedclusterinstances` Documentation

The `provisionedclusterinstances` SDK allows for interaction with Azure Resource Manager `hybridazurekubernetesservice` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/hybridazurekubernetesservice/2024-01-01/provisionedclusterinstances"
```


### Client Initialization

```go
client := provisionedclusterinstances.NewProvisionedClusterInstancesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProvisionedClusterInstancesClient.AgentPoolCreateOrUpdate`

```go
ctx := context.TODO()
id := provisionedclusterinstances.NewScopedAgentPoolID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "agentPoolName")

payload := provisionedclusterinstances.AgentPool{
	// ...
}


if err := client.AgentPoolCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ProvisionedClusterInstancesClient.AgentPoolDelete`

```go
ctx := context.TODO()
id := provisionedclusterinstances.NewScopedAgentPoolID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "agentPoolName")

if err := client.AgentPoolDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProvisionedClusterInstancesClient.AgentPoolGet`

```go
ctx := context.TODO()
id := provisionedclusterinstances.NewScopedAgentPoolID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "agentPoolName")

read, err := client.AgentPoolGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProvisionedClusterInstancesClient.AgentPoolListByProvisionedCluster`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.AgentPoolListByProvisionedCluster(ctx, id)` can be used to do batched pagination
items, err := client.AgentPoolListByProvisionedClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProvisionedClusterInstancesClient.DeleteKubernetesVersions`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

if err := client.DeleteKubernetesVersionsThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProvisionedClusterInstancesClient.DeleteVMSkus`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

if err := client.DeleteVMSkusThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProvisionedClusterInstancesClient.GetKubernetesVersions`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.GetKubernetesVersions(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProvisionedClusterInstancesClient.GetUpgradeProfile`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.GetUpgradeProfile(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProvisionedClusterInstancesClient.GetVMSkus`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.GetVMSkus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProvisionedClusterInstancesClient.HybridIdentityMetadataDelete`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

if err := client.HybridIdentityMetadataDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProvisionedClusterInstancesClient.HybridIdentityMetadataGet`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.HybridIdentityMetadataGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProvisionedClusterInstancesClient.HybridIdentityMetadataListByCluster`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.HybridIdentityMetadataListByCluster(ctx, id)` can be used to do batched pagination
items, err := client.HybridIdentityMetadataListByClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProvisionedClusterInstancesClient.HybridIdentityMetadataPut`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := provisionedclusterinstances.HybridIdentityMetadata{
	// ...
}


read, err := client.HybridIdentityMetadataPut(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProvisionedClusterInstancesClient.KubernetesVersionsList`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.KubernetesVersionsList(ctx, id)` can be used to do batched pagination
items, err := client.KubernetesVersionsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProvisionedClusterInstancesClient.ListAdminKubeconfig`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

if err := client.ListAdminKubeconfigThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProvisionedClusterInstancesClient.ListUserKubeconfig`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

if err := client.ListUserKubeconfigThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProvisionedClusterInstancesClient.ProvisionedClusterInstancesCreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := provisionedclusterinstances.ProvisionedCluster{
	// ...
}


if err := client.ProvisionedClusterInstancesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ProvisionedClusterInstancesClient.ProvisionedClusterInstancesDelete`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

if err := client.ProvisionedClusterInstancesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProvisionedClusterInstancesClient.ProvisionedClusterInstancesGet`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.ProvisionedClusterInstancesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProvisionedClusterInstancesClient.ProvisionedClusterInstancesList`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ProvisionedClusterInstancesList(ctx, id)` can be used to do batched pagination
items, err := client.ProvisionedClusterInstancesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProvisionedClusterInstancesClient.PutKubernetesVersions`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := provisionedclusterinstances.KubernetesVersionProfile{
	// ...
}


if err := client.PutKubernetesVersionsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ProvisionedClusterInstancesClient.PutVMSkus`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := provisionedclusterinstances.VMSkuProfile{
	// ...
}


if err := client.PutVMSkusThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ProvisionedClusterInstancesClient.VMSkusList`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.VMSkusList(ctx, id)` can be used to do batched pagination
items, err := client.VMSkusListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
