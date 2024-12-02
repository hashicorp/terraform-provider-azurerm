
## `github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-06-01-preview/nginxconfiguration` Documentation

The `nginxconfiguration` SDK allows for interaction with Azure Resource Manager `nginx` (API Version `2024-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-06-01-preview/nginxconfiguration"
```


### Client Initialization

```go
client := nginxconfiguration.NewNginxConfigurationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NginxConfigurationClient.ConfigurationsCreateOrUpdate`

```go
ctx := context.TODO()
id := nginxconfiguration.NewConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "configurationName")

payload := nginxconfiguration.NginxConfiguration{
	// ...
}


if err := client.ConfigurationsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NginxConfigurationClient.ConfigurationsDelete`

```go
ctx := context.TODO()
id := nginxconfiguration.NewConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "configurationName")

if err := client.ConfigurationsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NginxConfigurationClient.ConfigurationsGet`

```go
ctx := context.TODO()
id := nginxconfiguration.NewConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "configurationName")

read, err := client.ConfigurationsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NginxConfigurationClient.ConfigurationsList`

```go
ctx := context.TODO()
id := nginxconfiguration.NewNginxDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName")

// alternatively `client.ConfigurationsList(ctx, id)` can be used to do batched pagination
items, err := client.ConfigurationsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
