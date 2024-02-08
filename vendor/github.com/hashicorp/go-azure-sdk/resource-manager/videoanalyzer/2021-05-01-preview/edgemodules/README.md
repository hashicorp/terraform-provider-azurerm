
## `github.com/hashicorp/go-azure-sdk/resource-manager/videoanalyzer/2021-05-01-preview/edgemodules` Documentation

The `edgemodules` SDK allows for interaction with the Azure Resource Manager Service `videoanalyzer` (API Version `2021-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/videoanalyzer/2021-05-01-preview/edgemodules"
```


### Client Initialization

```go
client := edgemodules.NewEdgeModulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EdgeModulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := edgemodules.NewEdgeModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "videoAnalyzerValue", "edgeModuleValue")

payload := edgemodules.EdgeModuleEntity{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EdgeModulesClient.Delete`

```go
ctx := context.TODO()
id := edgemodules.NewEdgeModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "videoAnalyzerValue", "edgeModuleValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EdgeModulesClient.Get`

```go
ctx := context.TODO()
id := edgemodules.NewEdgeModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "videoAnalyzerValue", "edgeModuleValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EdgeModulesClient.List`

```go
ctx := context.TODO()
id := edgemodules.NewVideoAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "videoAnalyzerValue")

// alternatively `client.List(ctx, id, edgemodules.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, edgemodules.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EdgeModulesClient.ListProvisioningToken`

```go
ctx := context.TODO()
id := edgemodules.NewEdgeModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "videoAnalyzerValue", "edgeModuleValue")

payload := edgemodules.ListProvisioningTokenInput{
	// ...
}


read, err := client.ListProvisioningToken(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
