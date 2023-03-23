
## `github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01/nginxdeployment` Documentation

The `nginxdeployment` SDK allows for interaction with the Azure Resource Manager Service `nginx` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01/nginxdeployment"
```


### Client Initialization

```go
client := nginxdeployment.NewNginxDeploymentClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NginxDeploymentClient.DeploymentsCreateOrUpdate`

```go
ctx := context.TODO()
id := nginxdeployment.NewNginxDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentValue")

payload := nginxdeployment.NginxDeployment{
	// ...
}


if err := client.DeploymentsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NginxDeploymentClient.DeploymentsDelete`

```go
ctx := context.TODO()
id := nginxdeployment.NewNginxDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentValue")

if err := client.DeploymentsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NginxDeploymentClient.DeploymentsGet`

```go
ctx := context.TODO()
id := nginxdeployment.NewNginxDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentValue")

read, err := client.DeploymentsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NginxDeploymentClient.DeploymentsList`

```go
ctx := context.TODO()
id := nginxdeployment.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.DeploymentsList(ctx, id)` can be used to do batched pagination
items, err := client.DeploymentsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NginxDeploymentClient.DeploymentsListByResourceGroup`

```go
ctx := context.TODO()
id := nginxdeployment.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.DeploymentsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.DeploymentsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NginxDeploymentClient.DeploymentsUpdate`

```go
ctx := context.TODO()
id := nginxdeployment.NewNginxDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentValue")

payload := nginxdeployment.NginxDeploymentUpdateParameters{
	// ...
}


if err := client.DeploymentsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
