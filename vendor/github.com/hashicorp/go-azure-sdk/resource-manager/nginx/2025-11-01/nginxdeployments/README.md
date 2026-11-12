
## `github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2025-11-01/nginxdeployments` Documentation

The `nginxdeployments` SDK allows for interaction with Azure Resource Manager `nginx` (API Version `2025-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2025-11-01/nginxdeployments"
```


### Client Initialization

```go
client := nginxdeployments.NewNginxDeploymentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NginxDeploymentsClient.DefaultWafPolicyList`

```go
ctx := context.TODO()
id := nginxdeployments.NewNginxDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName")

// alternatively `client.DefaultWafPolicyList(ctx, id)` can be used to do batched pagination
items, err := client.DefaultWafPolicyListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NginxDeploymentsClient.DeploymentsCreateOrUpdate`

```go
ctx := context.TODO()
id := nginxdeployments.NewNginxDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName")

payload := nginxdeployments.NginxDeployment{
	// ...
}


if err := client.DeploymentsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NginxDeploymentsClient.DeploymentsDelete`

```go
ctx := context.TODO()
id := nginxdeployments.NewNginxDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName")

if err := client.DeploymentsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NginxDeploymentsClient.DeploymentsGet`

```go
ctx := context.TODO()
id := nginxdeployments.NewNginxDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName")

read, err := client.DeploymentsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NginxDeploymentsClient.DeploymentsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.DeploymentsList(ctx, id)` can be used to do batched pagination
items, err := client.DeploymentsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NginxDeploymentsClient.DeploymentsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.DeploymentsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.DeploymentsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NginxDeploymentsClient.DeploymentsUpdate`

```go
ctx := context.TODO()
id := nginxdeployments.NewNginxDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName")

payload := nginxdeployments.NginxDeploymentUpdateParameters{
	// ...
}


if err := client.DeploymentsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NginxDeploymentsClient.WafPolicyList`

```go
ctx := context.TODO()
id := nginxdeployments.NewNginxDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName")

// alternatively `client.WafPolicyList(ctx, id)` can be used to do batched pagination
items, err := client.WafPolicyListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
