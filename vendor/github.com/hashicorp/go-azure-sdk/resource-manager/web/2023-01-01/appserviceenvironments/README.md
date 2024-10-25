
## `github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/appserviceenvironments` Documentation

The `appserviceenvironments` SDK allows for interaction with Azure Resource Manager `web` (API Version `2023-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/appserviceenvironments"
```


### Client Initialization

```go
client := appserviceenvironments.NewAppServiceEnvironmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AppServiceEnvironmentsClient.ApproveOrRejectPrivateEndpointConnection`

```go
ctx := context.TODO()
id := appserviceenvironments.NewHostingEnvironmentPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName", "privateEndpointConnectionName")

payload := appserviceenvironments.RemotePrivateEndpointConnectionARMResource{
	// ...
}


if err := client.ApproveOrRejectPrivateEndpointConnectionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppServiceEnvironmentsClient.ChangeVnet`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

payload := appserviceenvironments.VirtualNetworkProfile{
	// ...
}


// alternatively `client.ChangeVnet(ctx, id, payload)` can be used to do batched pagination
items, err := client.ChangeVnetComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

payload := appserviceenvironments.AppServiceEnvironmentResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppServiceEnvironmentsClient.CreateOrUpdateMultiRolePool`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

payload := appserviceenvironments.WorkerPoolResource{
	// ...
}


if err := client.CreateOrUpdateMultiRolePoolThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppServiceEnvironmentsClient.CreateOrUpdateWorkerPool`

```go
ctx := context.TODO()
id := appserviceenvironments.NewWorkerPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName", "workerPoolName")

payload := appserviceenvironments.WorkerPoolResource{
	// ...
}


if err := client.CreateOrUpdateWorkerPoolThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppServiceEnvironmentsClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

if err := client.DeleteThenPoll(ctx, id, appserviceenvironments.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `AppServiceEnvironmentsClient.DeleteAseCustomDnsSuffixConfiguration`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

read, err := client.DeleteAseCustomDnsSuffixConfiguration(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.DeletePrivateEndpointConnection`

```go
ctx := context.TODO()
id := appserviceenvironments.NewHostingEnvironmentPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName", "privateEndpointConnectionName")

if err := client.DeletePrivateEndpointConnectionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppServiceEnvironmentsClient.Get`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.GetAseCustomDnsSuffixConfiguration`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

read, err := client.GetAseCustomDnsSuffixConfiguration(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.GetAseV3NetworkingConfiguration`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

read, err := client.GetAseV3NetworkingConfiguration(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.GetDiagnosticsItem`

```go
ctx := context.TODO()
id := appserviceenvironments.NewHostingEnvironmentDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName", "diagnosticName")

read, err := client.GetDiagnosticsItem(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.GetInboundNetworkDependenciesEndpoints`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

// alternatively `client.GetInboundNetworkDependenciesEndpoints(ctx, id)` can be used to do batched pagination
items, err := client.GetInboundNetworkDependenciesEndpointsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.GetMultiRolePool`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

read, err := client.GetMultiRolePool(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.GetOutboundNetworkDependenciesEndpoints`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

// alternatively `client.GetOutboundNetworkDependenciesEndpoints(ctx, id)` can be used to do batched pagination
items, err := client.GetOutboundNetworkDependenciesEndpointsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.GetPrivateEndpointConnection`

```go
ctx := context.TODO()
id := appserviceenvironments.NewHostingEnvironmentPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName", "privateEndpointConnectionName")

read, err := client.GetPrivateEndpointConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.GetPrivateEndpointConnectionList`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

// alternatively `client.GetPrivateEndpointConnectionList(ctx, id)` can be used to do batched pagination
items, err := client.GetPrivateEndpointConnectionListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.GetPrivateLinkResources`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

read, err := client.GetPrivateLinkResources(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.GetVipInfo`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

read, err := client.GetVipInfo(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.GetWorkerPool`

```go
ctx := context.TODO()
id := appserviceenvironments.NewWorkerPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName", "workerPoolName")

read, err := client.GetWorkerPool(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.List`

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


### Example Usage: `AppServiceEnvironmentsClient.ListAppServicePlans`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

// alternatively `client.ListAppServicePlans(ctx, id)` can be used to do batched pagination
items, err := client.ListAppServicePlansComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.ListByResourceGroup`

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


### Example Usage: `AppServiceEnvironmentsClient.ListCapacities`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

// alternatively `client.ListCapacities(ctx, id)` can be used to do batched pagination
items, err := client.ListCapacitiesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.ListDiagnostics`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

read, err := client.ListDiagnostics(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.ListMultiRoleMetricDefinitions`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

// alternatively `client.ListMultiRoleMetricDefinitions(ctx, id)` can be used to do batched pagination
items, err := client.ListMultiRoleMetricDefinitionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.ListMultiRolePoolInstanceMetricDefinitions`

```go
ctx := context.TODO()
id := appserviceenvironments.NewDefaultInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName", "instanceName")

// alternatively `client.ListMultiRolePoolInstanceMetricDefinitions(ctx, id)` can be used to do batched pagination
items, err := client.ListMultiRolePoolInstanceMetricDefinitionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.ListMultiRolePoolSkus`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

// alternatively `client.ListMultiRolePoolSkus(ctx, id)` can be used to do batched pagination
items, err := client.ListMultiRolePoolSkusComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.ListMultiRolePools`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

// alternatively `client.ListMultiRolePools(ctx, id)` can be used to do batched pagination
items, err := client.ListMultiRolePoolsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.ListMultiRoleUsages`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

// alternatively `client.ListMultiRoleUsages(ctx, id)` can be used to do batched pagination
items, err := client.ListMultiRoleUsagesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.ListOperations`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

read, err := client.ListOperations(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.ListUsages`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

// alternatively `client.ListUsages(ctx, id, appserviceenvironments.DefaultListUsagesOperationOptions())` can be used to do batched pagination
items, err := client.ListUsagesComplete(ctx, id, appserviceenvironments.DefaultListUsagesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.ListWebApps`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

// alternatively `client.ListWebApps(ctx, id, appserviceenvironments.DefaultListWebAppsOperationOptions())` can be used to do batched pagination
items, err := client.ListWebAppsComplete(ctx, id, appserviceenvironments.DefaultListWebAppsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.ListWebWorkerMetricDefinitions`

```go
ctx := context.TODO()
id := appserviceenvironments.NewWorkerPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName", "workerPoolName")

// alternatively `client.ListWebWorkerMetricDefinitions(ctx, id)` can be used to do batched pagination
items, err := client.ListWebWorkerMetricDefinitionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.ListWebWorkerUsages`

```go
ctx := context.TODO()
id := appserviceenvironments.NewWorkerPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName", "workerPoolName")

// alternatively `client.ListWebWorkerUsages(ctx, id)` can be used to do batched pagination
items, err := client.ListWebWorkerUsagesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.ListWorkerPoolInstanceMetricDefinitions`

```go
ctx := context.TODO()
id := appserviceenvironments.NewWorkerPoolInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName", "workerPoolName", "instanceName")

// alternatively `client.ListWorkerPoolInstanceMetricDefinitions(ctx, id)` can be used to do batched pagination
items, err := client.ListWorkerPoolInstanceMetricDefinitionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.ListWorkerPoolSkus`

```go
ctx := context.TODO()
id := appserviceenvironments.NewWorkerPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName", "workerPoolName")

// alternatively `client.ListWorkerPoolSkus(ctx, id)` can be used to do batched pagination
items, err := client.ListWorkerPoolSkusComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.ListWorkerPools`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

// alternatively `client.ListWorkerPools(ctx, id)` can be used to do batched pagination
items, err := client.ListWorkerPoolsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.Reboot`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

read, err := client.Reboot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.Resume`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

// alternatively `client.Resume(ctx, id)` can be used to do batched pagination
items, err := client.ResumeComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.Suspend`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

// alternatively `client.Suspend(ctx, id)` can be used to do batched pagination
items, err := client.SuspendComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceEnvironmentsClient.TestUpgradeAvailableNotification`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

read, err := client.TestUpgradeAvailableNotification(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.Update`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

payload := appserviceenvironments.AppServiceEnvironmentPatchResource{
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


### Example Usage: `AppServiceEnvironmentsClient.UpdateAseCustomDnsSuffixConfiguration`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

payload := appserviceenvironments.CustomDnsSuffixConfiguration{
	// ...
}


read, err := client.UpdateAseCustomDnsSuffixConfiguration(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.UpdateAseNetworkingConfiguration`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

payload := appserviceenvironments.AseV3NetworkingConfiguration{
	// ...
}


read, err := client.UpdateAseNetworkingConfiguration(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.UpdateMultiRolePool`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

payload := appserviceenvironments.WorkerPoolResource{
	// ...
}


read, err := client.UpdateMultiRolePool(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.UpdateWorkerPool`

```go
ctx := context.TODO()
id := appserviceenvironments.NewWorkerPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName", "workerPoolName")

payload := appserviceenvironments.WorkerPoolResource{
	// ...
}


read, err := client.UpdateWorkerPool(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceEnvironmentsClient.Upgrade`

```go
ctx := context.TODO()
id := commonids.NewAppServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostingEnvironmentName")

if err := client.UpgradeThenPoll(ctx, id); err != nil {
	// handle the error
}
```
