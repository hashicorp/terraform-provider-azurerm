
## `github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-06-01-preview/nginxconfigurationanalysis` Documentation

The `nginxconfigurationanalysis` SDK allows for interaction with Azure Resource Manager `nginx` (API Version `2024-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-06-01-preview/nginxconfigurationanalysis"
```


### Client Initialization

```go
client := nginxconfigurationanalysis.NewNginxConfigurationAnalysisClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NginxConfigurationAnalysisClient.ConfigurationsAnalysis`

```go
ctx := context.TODO()
id := nginxconfigurationanalysis.NewConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "configurationName")

payload := nginxconfigurationanalysis.AnalysisCreate{
	// ...
}


read, err := client.ConfigurationsAnalysis(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
