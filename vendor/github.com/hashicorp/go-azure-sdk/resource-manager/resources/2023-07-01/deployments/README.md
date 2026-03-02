
## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/deployments` Documentation

The `deployments` SDK allows for interaction with Azure Resource Manager `resources` (API Version `2023-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/deployments"
```


### Client Initialization

```go
client := deployments.NewDeploymentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeploymentsClient.CalculateTemplateHash`

```go
ctx := context.TODO()
var payload interface{}

read, err := client.CalculateTemplateHash(ctx, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.Cancel`

```go
ctx := context.TODO()
id := deployments.NewResourceGroupProviderDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentName")

read, err := client.Cancel(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.CancelAtManagementGroupScope`

```go
ctx := context.TODO()
id := deployments.NewProviders2DeploymentID("groupId", "deploymentName")

read, err := client.CancelAtManagementGroupScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.CancelAtScope`

```go
ctx := context.TODO()
id := deployments.NewScopedDeploymentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "deploymentName")

read, err := client.CancelAtScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.CancelAtSubscriptionScope`

```go
ctx := context.TODO()
id := deployments.NewProviderDeploymentID("12345678-1234-9876-4563-123456789012", "deploymentName")

read, err := client.CancelAtSubscriptionScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.CancelAtTenantScope`

```go
ctx := context.TODO()
id := deployments.NewDeploymentID("deploymentName")

read, err := client.CancelAtTenantScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.CheckExistence`

```go
ctx := context.TODO()
id := deployments.NewResourceGroupProviderDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentName")

read, err := client.CheckExistence(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.CheckExistenceAtManagementGroupScope`

```go
ctx := context.TODO()
id := deployments.NewProviders2DeploymentID("groupId", "deploymentName")

read, err := client.CheckExistenceAtManagementGroupScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.CheckExistenceAtScope`

```go
ctx := context.TODO()
id := deployments.NewScopedDeploymentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "deploymentName")

read, err := client.CheckExistenceAtScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.CheckExistenceAtSubscriptionScope`

```go
ctx := context.TODO()
id := deployments.NewProviderDeploymentID("12345678-1234-9876-4563-123456789012", "deploymentName")

read, err := client.CheckExistenceAtSubscriptionScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.CheckExistenceAtTenantScope`

```go
ctx := context.TODO()
id := deployments.NewDeploymentID("deploymentName")

read, err := client.CheckExistenceAtTenantScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := deployments.NewResourceGroupProviderDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentName")

payload := deployments.Deployment{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.CreateOrUpdateAtManagementGroupScope`

```go
ctx := context.TODO()
id := deployments.NewProviders2DeploymentID("groupId", "deploymentName")

payload := deployments.ScopedDeployment{
	// ...
}


if err := client.CreateOrUpdateAtManagementGroupScopeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.CreateOrUpdateAtScope`

```go
ctx := context.TODO()
id := deployments.NewScopedDeploymentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "deploymentName")

payload := deployments.Deployment{
	// ...
}


if err := client.CreateOrUpdateAtScopeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.CreateOrUpdateAtSubscriptionScope`

```go
ctx := context.TODO()
id := deployments.NewProviderDeploymentID("12345678-1234-9876-4563-123456789012", "deploymentName")

payload := deployments.Deployment{
	// ...
}


if err := client.CreateOrUpdateAtSubscriptionScopeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.CreateOrUpdateAtTenantScope`

```go
ctx := context.TODO()
id := deployments.NewDeploymentID("deploymentName")

payload := deployments.ScopedDeployment{
	// ...
}


if err := client.CreateOrUpdateAtTenantScopeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.Delete`

```go
ctx := context.TODO()
id := deployments.NewResourceGroupProviderDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.DeleteAtManagementGroupScope`

```go
ctx := context.TODO()
id := deployments.NewProviders2DeploymentID("groupId", "deploymentName")

if err := client.DeleteAtManagementGroupScopeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.DeleteAtScope`

```go
ctx := context.TODO()
id := deployments.NewScopedDeploymentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "deploymentName")

if err := client.DeleteAtScopeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.DeleteAtSubscriptionScope`

```go
ctx := context.TODO()
id := deployments.NewProviderDeploymentID("12345678-1234-9876-4563-123456789012", "deploymentName")

if err := client.DeleteAtSubscriptionScopeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.DeleteAtTenantScope`

```go
ctx := context.TODO()
id := deployments.NewDeploymentID("deploymentName")

if err := client.DeleteAtTenantScopeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.ExportTemplate`

```go
ctx := context.TODO()
id := deployments.NewResourceGroupProviderDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentName")

read, err := client.ExportTemplate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.ExportTemplateAtManagementGroupScope`

```go
ctx := context.TODO()
id := deployments.NewProviders2DeploymentID("groupId", "deploymentName")

read, err := client.ExportTemplateAtManagementGroupScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.ExportTemplateAtScope`

```go
ctx := context.TODO()
id := deployments.NewScopedDeploymentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "deploymentName")

read, err := client.ExportTemplateAtScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.ExportTemplateAtSubscriptionScope`

```go
ctx := context.TODO()
id := deployments.NewProviderDeploymentID("12345678-1234-9876-4563-123456789012", "deploymentName")

read, err := client.ExportTemplateAtSubscriptionScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.ExportTemplateAtTenantScope`

```go
ctx := context.TODO()
id := deployments.NewDeploymentID("deploymentName")

read, err := client.ExportTemplateAtTenantScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.Get`

```go
ctx := context.TODO()
id := deployments.NewResourceGroupProviderDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.GetAtManagementGroupScope`

```go
ctx := context.TODO()
id := deployments.NewProviders2DeploymentID("groupId", "deploymentName")

read, err := client.GetAtManagementGroupScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.GetAtScope`

```go
ctx := context.TODO()
id := deployments.NewScopedDeploymentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "deploymentName")

read, err := client.GetAtScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.GetAtSubscriptionScope`

```go
ctx := context.TODO()
id := deployments.NewProviderDeploymentID("12345678-1234-9876-4563-123456789012", "deploymentName")

read, err := client.GetAtSubscriptionScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.GetAtTenantScope`

```go
ctx := context.TODO()
id := deployments.NewDeploymentID("deploymentName")

read, err := client.GetAtTenantScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.ListAtManagementGroupScope`

```go
ctx := context.TODO()
id := commonids.NewManagementGroupID("groupId")

// alternatively `client.ListAtManagementGroupScope(ctx, id, deployments.DefaultListAtManagementGroupScopeOperationOptions())` can be used to do batched pagination
items, err := client.ListAtManagementGroupScopeComplete(ctx, id, deployments.DefaultListAtManagementGroupScopeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeploymentsClient.ListAtScope`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ListAtScope(ctx, id, deployments.DefaultListAtScopeOperationOptions())` can be used to do batched pagination
items, err := client.ListAtScopeComplete(ctx, id, deployments.DefaultListAtScopeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeploymentsClient.ListAtSubscriptionScope`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAtSubscriptionScope(ctx, id, deployments.DefaultListAtSubscriptionScopeOperationOptions())` can be used to do batched pagination
items, err := client.ListAtSubscriptionScopeComplete(ctx, id, deployments.DefaultListAtSubscriptionScopeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeploymentsClient.ListAtTenantScope`

```go
ctx := context.TODO()


// alternatively `client.ListAtTenantScope(ctx, deployments.DefaultListAtTenantScopeOperationOptions())` can be used to do batched pagination
items, err := client.ListAtTenantScopeComplete(ctx, deployments.DefaultListAtTenantScopeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeploymentsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, deployments.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, deployments.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeploymentsClient.Validate`

```go
ctx := context.TODO()
id := deployments.NewResourceGroupProviderDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentName")

payload := deployments.Deployment{
	// ...
}


if err := client.ValidateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.ValidateAtManagementGroupScope`

```go
ctx := context.TODO()
id := deployments.NewProviders2DeploymentID("groupId", "deploymentName")

payload := deployments.ScopedDeployment{
	// ...
}


if err := client.ValidateAtManagementGroupScopeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.ValidateAtScope`

```go
ctx := context.TODO()
id := deployments.NewScopedDeploymentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "deploymentName")

payload := deployments.Deployment{
	// ...
}


if err := client.ValidateAtScopeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.ValidateAtSubscriptionScope`

```go
ctx := context.TODO()
id := deployments.NewProviderDeploymentID("12345678-1234-9876-4563-123456789012", "deploymentName")

payload := deployments.Deployment{
	// ...
}


if err := client.ValidateAtSubscriptionScopeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.ValidateAtTenantScope`

```go
ctx := context.TODO()
id := deployments.NewDeploymentID("deploymentName")

payload := deployments.ScopedDeployment{
	// ...
}


if err := client.ValidateAtTenantScopeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.WhatIf`

```go
ctx := context.TODO()
id := deployments.NewResourceGroupProviderDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentName")

payload := deployments.DeploymentWhatIf{
	// ...
}


if err := client.WhatIfThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.WhatIfAtManagementGroupScope`

```go
ctx := context.TODO()
id := deployments.NewProviders2DeploymentID("groupId", "deploymentName")

payload := deployments.ScopedDeploymentWhatIf{
	// ...
}


if err := client.WhatIfAtManagementGroupScopeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.WhatIfAtSubscriptionScope`

```go
ctx := context.TODO()
id := deployments.NewProviderDeploymentID("12345678-1234-9876-4563-123456789012", "deploymentName")

payload := deployments.DeploymentWhatIf{
	// ...
}


if err := client.WhatIfAtSubscriptionScopeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.WhatIfAtTenantScope`

```go
ctx := context.TODO()
id := deployments.NewDeploymentID("deploymentName")

payload := deployments.ScopedDeploymentWhatIf{
	// ...
}


if err := client.WhatIfAtTenantScopeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
