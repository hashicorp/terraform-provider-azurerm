
## `github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2023-05-01/deployments` Documentation

The `deployments` SDK allows for interaction with Azure Resource Manager `cognitive` (API Version `2023-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2023-05-01/deployments"
```


### Client Initialization

```go
client := deployments.NewDeploymentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeploymentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := deployments.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "deploymentName")

payload := deployments.Deployment{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.Delete`

```go
ctx := context.TODO()
id := deployments.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "deploymentName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentsClient.Get`

```go
ctx := context.TODO()
id := deployments.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "deploymentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentsClient.List`

```go
ctx := context.TODO()
id := deployments.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
