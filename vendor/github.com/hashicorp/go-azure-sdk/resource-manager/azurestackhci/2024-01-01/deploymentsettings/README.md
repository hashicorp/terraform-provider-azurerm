
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/deploymentsettings` Documentation

The `deploymentsettings` SDK allows for interaction with Azure Resource Manager `azurestackhci` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/deploymentsettings"
```


### Client Initialization

```go
client := deploymentsettings.NewDeploymentSettingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeploymentSettingsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := deploymentsettings.NewDeploymentSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "deploymentSettingName")

payload := deploymentsettings.DeploymentSetting{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentSettingsClient.Delete`

```go
ctx := context.TODO()
id := deploymentsettings.NewDeploymentSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "deploymentSettingName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DeploymentSettingsClient.Get`

```go
ctx := context.TODO()
id := deploymentsettings.NewDeploymentSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "deploymentSettingName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeploymentSettingsClient.ListByClusters`

```go
ctx := context.TODO()
id := deploymentsettings.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

// alternatively `client.ListByClusters(ctx, id)` can be used to do batched pagination
items, err := client.ListByClustersComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
