
## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-02-01/templatespecversions` Documentation

The `templatespecversions` SDK allows for interaction with Azure Resource Manager `resources` (API Version `2022-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-02-01/templatespecversions"
```


### Client Initialization

```go
client := templatespecversions.NewTemplateSpecVersionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TemplateSpecVersionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := templatespecversions.NewTemplateSpecVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "templateSpecName", "versionName")

payload := templatespecversions.TemplateSpecVersion{
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


### Example Usage: `TemplateSpecVersionsClient.Delete`

```go
ctx := context.TODO()
id := templatespecversions.NewTemplateSpecVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "templateSpecName", "versionName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TemplateSpecVersionsClient.Get`

```go
ctx := context.TODO()
id := templatespecversions.NewTemplateSpecVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "templateSpecName", "versionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TemplateSpecVersionsClient.GetBuiltIn`

```go
ctx := context.TODO()
id := templatespecversions.NewVersionID("builtInTemplateSpecName", "versionName")

read, err := client.GetBuiltIn(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TemplateSpecVersionsClient.List`

```go
ctx := context.TODO()
id := templatespecversions.NewTemplateSpecID("12345678-1234-9876-4563-123456789012", "example-resource-group", "templateSpecName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TemplateSpecVersionsClient.ListBuiltIns`

```go
ctx := context.TODO()
id := templatespecversions.NewBuiltInTemplateSpecID("builtInTemplateSpecName")

// alternatively `client.ListBuiltIns(ctx, id)` can be used to do batched pagination
items, err := client.ListBuiltInsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TemplateSpecVersionsClient.Update`

```go
ctx := context.TODO()
id := templatespecversions.NewTemplateSpecVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "templateSpecName", "versionName")

payload := templatespecversions.TemplateSpecVersionUpdateModel{
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
