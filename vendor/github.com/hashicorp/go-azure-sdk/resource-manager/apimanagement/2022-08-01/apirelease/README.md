
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apirelease` Documentation

The `apirelease` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apirelease"
```


### Client Initialization

```go
client := apirelease.NewApiReleaseClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiReleaseClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apirelease.NewReleaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "releaseId")

payload := apirelease.ApiReleaseContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, apirelease.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiReleaseClient.Delete`

```go
ctx := context.TODO()
id := apirelease.NewReleaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "releaseId")

read, err := client.Delete(ctx, id, apirelease.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiReleaseClient.Get`

```go
ctx := context.TODO()
id := apirelease.NewReleaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "releaseId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiReleaseClient.GetEntityTag`

```go
ctx := context.TODO()
id := apirelease.NewReleaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "releaseId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiReleaseClient.ListByService`

```go
ctx := context.TODO()
id := apirelease.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

// alternatively `client.ListByService(ctx, id, apirelease.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, apirelease.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiReleaseClient.Update`

```go
ctx := context.TODO()
id := apirelease.NewReleaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "releaseId")

payload := apirelease.ApiReleaseContract{
	// ...
}


read, err := client.Update(ctx, id, payload, apirelease.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
