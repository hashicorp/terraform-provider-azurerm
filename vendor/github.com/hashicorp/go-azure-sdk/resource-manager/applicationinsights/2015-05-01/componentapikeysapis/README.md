
## `github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentapikeysapis` Documentation

The `componentapikeysapis` SDK allows for interaction with Azure Resource Manager `applicationinsights` (API Version `2015-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentapikeysapis"
```


### Client Initialization

```go
client := componentapikeysapis.NewComponentApiKeysAPIsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ComponentApiKeysAPIsClient.APIKeysCreate`

```go
ctx := context.TODO()
id := componentapikeysapis.NewComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName")

payload := componentapikeysapis.APIKeyRequest{
	// ...
}


read, err := client.APIKeysCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentApiKeysAPIsClient.APIKeysDelete`

```go
ctx := context.TODO()
id := componentapikeysapis.NewApiKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName", "keyId")

read, err := client.APIKeysDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentApiKeysAPIsClient.APIKeysGet`

```go
ctx := context.TODO()
id := componentapikeysapis.NewApiKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName", "keyId")

read, err := client.APIKeysGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentApiKeysAPIsClient.APIKeysList`

```go
ctx := context.TODO()
id := componentapikeysapis.NewComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName")

read, err := client.APIKeysList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
