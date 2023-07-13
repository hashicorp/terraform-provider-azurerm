
## `github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/metadata` Documentation

The `metadata` SDK allows for interaction with the Azure Resource Manager Service `securityinsights` (API Version `2022-10-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/metadata"
```


### Client Initialization

```go
client := metadata.NewMetadataClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MetadataClient.Create`

```go
ctx := context.TODO()
id := metadata.NewMetadataID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "metadataValue")

payload := metadata.MetadataModel{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MetadataClient.Delete`

```go
ctx := context.TODO()
id := metadata.NewMetadataID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "metadataValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MetadataClient.Get`

```go
ctx := context.TODO()
id := metadata.NewMetadataID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "metadataValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MetadataClient.List`

```go
ctx := context.TODO()
id := metadata.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.List(ctx, id, metadata.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, metadata.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MetadataClient.Update`

```go
ctx := context.TODO()
id := metadata.NewMetadataID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "metadataValue")

payload := metadata.MetadataPatch{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
