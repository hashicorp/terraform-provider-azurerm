
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiwiki` Documentation

The `apiwiki` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiwiki"
```


### Client Initialization

```go
client := apiwiki.NewApiWikiClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiWikiClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apiwiki.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

payload := apiwiki.WikiContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, apiwiki.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiWikiClient.Delete`

```go
ctx := context.TODO()
id := apiwiki.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

read, err := client.Delete(ctx, id, apiwiki.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiWikiClient.Get`

```go
ctx := context.TODO()
id := apiwiki.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiWikiClient.GetEntityTag`

```go
ctx := context.TODO()
id := apiwiki.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiWikiClient.List`

```go
ctx := context.TODO()
id := apiwiki.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

// alternatively `client.List(ctx, id, apiwiki.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, apiwiki.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiWikiClient.Update`

```go
ctx := context.TODO()
id := apiwiki.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

payload := apiwiki.WikiUpdateContract{
	// ...
}


read, err := client.Update(ctx, id, payload, apiwiki.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
