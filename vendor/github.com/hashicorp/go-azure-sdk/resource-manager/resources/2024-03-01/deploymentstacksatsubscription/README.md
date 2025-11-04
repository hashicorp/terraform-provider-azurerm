
## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2024-03-01/deploymentstacksatsubscription` Documentation

The `deploymentstacksatsubscription` SDK allows for interaction with Azure Resource Manager `resources` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2024-03-01/deploymentstacksatsubscription"
```


### Client Initialization

```go
client := deploymentstacksatsubscription.NewDeploymentStacksAtSubscriptionClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeploymentStacksAtSubscriptionClient.DeploymentStacksCreateOrUpdateAtSubscription`

```go
ctx := context.TODO()
id := deploymentstacksatsubscription.NewDeploymentStackID("12345678-1234-9876-4563-123456789012", "deploymentStackName")

payload := deploymentstacksatsubscription.DeploymentStack{
	// ...
}


if err := client.DeploymentStacksCreateOrUpdateAtSubscriptionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentStacksAtSubscriptionClient.DeploymentStacksDeleteAtSubscription`

```go
ctx := context.TODO()
id := deploymentstacksatsubscription.NewDeploymentStackID("12345678-1234-9876-4563-123456789012", "deploymentStackName")

if err := client.DeploymentStacksDeleteAtSubscriptionThenPoll(ctx, id, deploymentstacksatsubscription.DefaultDeploymentStacksDeleteAtSubscriptionOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentStacksAtSubscriptionClient.DeploymentStacksExportTemplateAtSubscription`

```go
ctx := context.TODO()
id := deploymentstacksatsubscription.NewDeploymentStackID("12345678-1234-9876-4563-123456789012", "deploymentStackName")

read, err := client.DeploymentStacksExportTemplateAtSubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentStacksAtSubscriptionClient.DeploymentStacksGetAtSubscription`

```go
ctx := context.TODO()
id := deploymentstacksatsubscription.NewDeploymentStackID("12345678-1234-9876-4563-123456789012", "deploymentStackName")

read, err := client.DeploymentStacksGetAtSubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentStacksAtSubscriptionClient.DeploymentStacksListAtSubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.DeploymentStacksListAtSubscription(ctx, id)` can be used to do batched pagination
items, err := client.DeploymentStacksListAtSubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeploymentStacksAtSubscriptionClient.DeploymentStacksValidateStackAtSubscription`

```go
ctx := context.TODO()
id := deploymentstacksatsubscription.NewDeploymentStackID("12345678-1234-9876-4563-123456789012", "deploymentStackName")

payload := deploymentstacksatsubscription.DeploymentStack{
	// ...
}


if err := client.DeploymentStacksValidateStackAtSubscriptionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
