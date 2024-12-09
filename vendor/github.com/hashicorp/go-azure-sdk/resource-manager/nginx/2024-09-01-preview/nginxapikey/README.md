
## `github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-09-01-preview/nginxapikey` Documentation

The `nginxapikey` SDK allows for interaction with Azure Resource Manager `nginx` (API Version `2024-09-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-09-01-preview/nginxapikey"
```


### Client Initialization

```go
client := nginxapikey.NewNginxApiKeyClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NginxApiKeyClient.ApiKeysCreateOrUpdate`

```go
ctx := context.TODO()
id := nginxapikey.NewApiKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "apiKeyName")

payload := nginxapikey.NginxDeploymentApiKeyRequest{
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


### Example Usage: `NginxApiKeyClient.ApiKeysDelete`

```go
ctx := context.TODO()
id := nginxapikey.NewApiKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "apiKeyName")

read, err := client.ApiKeysDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NginxApiKeyClient.ApiKeysGet`

```go
ctx := context.TODO()
id := nginxapikey.NewApiKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "apiKeyName")

read, err := client.ApiKeysGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NginxApiKeyClient.ApiKeysList`

```go
ctx := context.TODO()
id := nginxapikey.NewNginxDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName")

// alternatively `client.ApiKeysList(ctx, id)` can be used to do batched pagination
items, err := client.ApiKeysListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
