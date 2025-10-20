
## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2024-03-01/deploymentstacksatmanagementgroup` Documentation

The `deploymentstacksatmanagementgroup` SDK allows for interaction with Azure Resource Manager `resources` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2024-03-01/deploymentstacksatmanagementgroup"
```


### Client Initialization

```go
client := deploymentstacksatmanagementgroup.NewDeploymentStacksAtManagementGroupClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeploymentStacksAtManagementGroupClient.DeploymentStacksCreateOrUpdateAtManagementGroup`

```go
ctx := context.TODO()
id := deploymentstacksatmanagementgroup.NewProviders2DeploymentStackID("managementGroupId", "deploymentStackName")

payload := deploymentstacksatmanagementgroup.DeploymentStack{
	// ...
}


if err := client.DeploymentStacksCreateOrUpdateAtManagementGroupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentStacksAtManagementGroupClient.DeploymentStacksDeleteAtManagementGroup`

```go
ctx := context.TODO()
id := deploymentstacksatmanagementgroup.NewProviders2DeploymentStackID("managementGroupId", "deploymentStackName")

if err := client.DeploymentStacksDeleteAtManagementGroupThenPoll(ctx, id, deploymentstacksatmanagementgroup.DefaultDeploymentStacksDeleteAtManagementGroupOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentStacksAtManagementGroupClient.DeploymentStacksExportTemplateAtManagementGroup`

```go
ctx := context.TODO()
id := deploymentstacksatmanagementgroup.NewProviders2DeploymentStackID("managementGroupId", "deploymentStackName")

read, err := client.DeploymentStacksExportTemplateAtManagementGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentStacksAtManagementGroupClient.DeploymentStacksGetAtManagementGroup`

```go
ctx := context.TODO()
id := deploymentstacksatmanagementgroup.NewProviders2DeploymentStackID("managementGroupId", "deploymentStackName")

read, err := client.DeploymentStacksGetAtManagementGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentStacksAtManagementGroupClient.DeploymentStacksListAtManagementGroup`

```go
ctx := context.TODO()
id := commonids.NewManagementGroupID("groupId")

// alternatively `client.DeploymentStacksListAtManagementGroup(ctx, id)` can be used to do batched pagination
items, err := client.DeploymentStacksListAtManagementGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeploymentStacksAtManagementGroupClient.DeploymentStacksValidateStackAtManagementGroup`

```go
ctx := context.TODO()
id := deploymentstacksatmanagementgroup.NewProviders2DeploymentStackID("managementGroupId", "deploymentStackName")

payload := deploymentstacksatmanagementgroup.DeploymentStack{
	// ...
}


if err := client.DeploymentStacksValidateStackAtManagementGroupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
