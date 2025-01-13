
## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-10-01/deploymentscripts` Documentation

The `deploymentscripts` SDK allows for interaction with Azure Resource Manager `resources` (API Version `2020-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-10-01/deploymentscripts"
```


### Client Initialization

```go
client := deploymentscripts.NewDeploymentScriptsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeploymentScriptsClient.Create`

```go
ctx := context.TODO()
id := deploymentscripts.NewDeploymentScriptID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentScriptName")

payload := deploymentscripts.DeploymentScript{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentScriptsClient.Delete`

```go
ctx := context.TODO()
id := deploymentscripts.NewDeploymentScriptID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentScriptName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentScriptsClient.Get`

```go
ctx := context.TODO()
id := deploymentscripts.NewDeploymentScriptID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentScriptName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentScriptsClient.GetLogs`

```go
ctx := context.TODO()
id := deploymentscripts.NewDeploymentScriptID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentScriptName")

read, err := client.GetLogs(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentScriptsClient.GetLogsDefault`

```go
ctx := context.TODO()
id := deploymentscripts.NewDeploymentScriptID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentScriptName")

read, err := client.GetLogsDefault(ctx, id, deploymentscripts.DefaultGetLogsDefaultOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentScriptsClient.ListByResourceGroup`

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


### Example Usage: `DeploymentScriptsClient.ListBySubscription`

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


### Example Usage: `DeploymentScriptsClient.Update`

```go
ctx := context.TODO()
id := deploymentscripts.NewDeploymentScriptID("12345678-1234-9876-4563-123456789012", "example-resource-group", "deploymentScriptName")

payload := deploymentscripts.DeploymentScriptUpdateParameter{
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
