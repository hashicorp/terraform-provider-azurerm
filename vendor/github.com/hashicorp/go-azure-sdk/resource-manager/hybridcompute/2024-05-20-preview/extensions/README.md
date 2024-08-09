
## `github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-05-20-preview/extensions` Documentation

The `extensions` SDK allows for interaction with the Azure Resource Manager Service `hybridcompute` (API Version `2024-05-20-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-05-20-preview/extensions"
```


### Client Initialization

```go
client := extensions.NewExtensionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExtensionsClient.ExtensionMetadataGet`

```go
ctx := context.TODO()
id := extensions.NewVersionID("12345678-1234-9876-4563-123456789012", "locationValue", "publisherValue", "extensionTypeValue", "versionValue")

read, err := client.ExtensionMetadataGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExtensionsClient.ExtensionMetadataList`

```go
ctx := context.TODO()
id := extensions.NewExtensionTypeID("12345678-1234-9876-4563-123456789012", "locationValue", "publisherValue", "extensionTypeValue")

read, err := client.ExtensionMetadataList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
