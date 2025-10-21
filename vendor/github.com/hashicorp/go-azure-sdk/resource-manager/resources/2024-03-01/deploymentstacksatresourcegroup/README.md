
## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2024-03-01/deploymentstacksatresourcegroup` Documentation

The `deploymentstacksatresourcegroup` SDK allows for interaction with Azure Resource Manager `resources` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2024-03-01/deploymentstacksatresourcegroup"
```


### Client Initialization

```go
client := deploymentstacksatresourcegroup.NewDeploymentStacksAtResourceGroupClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeploymentStacksAtResourceGroupClient.DeploymentStacksCreateOrUpdateAtResourceGroup`

```go
ctx := context.TODO()
id := deploymentstacksatresourcegroup.NewProviderDeploymentStackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentStackName")

payload := deploymentstacksatresourcegroup.DeploymentStack{
	// ...
}


if err := client.DeploymentStacksCreateOrUpdateAtResourceGroupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentStacksAtResourceGroupClient.DeploymentStacksDeleteAtResourceGroup`

```go
ctx := context.TODO()
id := deploymentstacksatresourcegroup.NewProviderDeploymentStackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentStackName")

if err := client.DeploymentStacksDeleteAtResourceGroupThenPoll(ctx, id, deploymentstacksatresourcegroup.DefaultDeploymentStacksDeleteAtResourceGroupOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentStacksAtResourceGroupClient.DeploymentStacksExportTemplateAtResourceGroup`

```go
ctx := context.TODO()
id := deploymentstacksatresourcegroup.NewProviderDeploymentStackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentStackName")

read, err := client.DeploymentStacksExportTemplateAtResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentStacksAtResourceGroupClient.DeploymentStacksGetAtResourceGroup`

```go
ctx := context.TODO()
id := deploymentstacksatresourcegroup.NewProviderDeploymentStackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentStackName")

read, err := client.DeploymentStacksGetAtResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentStacksAtResourceGroupClient.DeploymentStacksListAtResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.DeploymentStacksListAtResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.DeploymentStacksListAtResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeploymentStacksAtResourceGroupClient.DeploymentStacksValidateStackAtResourceGroup`

```go
ctx := context.TODO()
id := deploymentstacksatresourcegroup.NewProviderDeploymentStackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentStackName")

payload := deploymentstacksatresourcegroup.DeploymentStack{
	// ...
}


if err := client.DeploymentStacksValidateStackAtResourceGroupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
