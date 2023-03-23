
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/fileservice` Documentation

The `fileservice` SDK allows for interaction with the Azure Resource Manager Service `storage` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/fileservice"
```


### Client Initialization

```go
client := fileservice.NewFileServiceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FileServiceClient.GetServiceProperties`

```go
ctx := context.TODO()
id := fileservice.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

read, err := client.GetServiceProperties(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FileServiceClient.List`

```go
ctx := context.TODO()
id := fileservice.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FileServiceClient.SetServiceProperties`

```go
ctx := context.TODO()
id := fileservice.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

payload := fileservice.FileServiceProperties{
	// ...
}


read, err := client.SetServiceProperties(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
