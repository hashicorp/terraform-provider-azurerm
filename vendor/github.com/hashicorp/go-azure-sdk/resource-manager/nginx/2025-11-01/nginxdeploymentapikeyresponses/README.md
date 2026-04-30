
## `github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2025-11-01/nginxdeploymentapikeyresponses` Documentation

The `nginxdeploymentapikeyresponses` SDK allows for interaction with Azure Resource Manager `nginx` (API Version `2025-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2025-11-01/nginxdeploymentapikeyresponses"
```


### Client Initialization

```go
client := nginxdeploymentapikeyresponses.NewNginxDeploymentApiKeyResponsesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NginxDeploymentApiKeyResponsesClient.ApiKeysCreateOrUpdate`

```go
ctx := context.TODO()
id := nginxdeploymentapikeyresponses.NewApiKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "apiKeyName")

payload := nginxdeploymentapikeyresponses.NginxDeploymentApiKeyRequest{
	// ...
}


read, err := client.ApiKeysCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NginxDeploymentApiKeyResponsesClient.ApiKeysDelete`

```go
ctx := context.TODO()
id := nginxdeploymentapikeyresponses.NewApiKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "apiKeyName")

read, err := client.ApiKeysDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NginxDeploymentApiKeyResponsesClient.ApiKeysGet`

```go
ctx := context.TODO()
id := nginxdeploymentapikeyresponses.NewApiKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "apiKeyName")

read, err := client.ApiKeysGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NginxDeploymentApiKeyResponsesClient.ApiKeysList`

```go
ctx := context.TODO()
id := nginxdeploymentapikeyresponses.NewNginxDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName")

// alternatively `client.ApiKeysList(ctx, id)` can be used to do batched pagination
items, err := client.ApiKeysListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
