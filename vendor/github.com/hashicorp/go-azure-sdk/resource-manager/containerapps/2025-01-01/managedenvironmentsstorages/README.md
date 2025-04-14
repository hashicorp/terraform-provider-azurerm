
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/managedenvironmentsstorages` Documentation

The `managedenvironmentsstorages` SDK allows for interaction with Azure Resource Manager `containerapps` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/managedenvironmentsstorages"
```


### Client Initialization

```go
client := managedenvironmentsstorages.NewManagedEnvironmentsStoragesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedEnvironmentsStoragesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := managedenvironmentsstorages.NewStorageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName", "storageName")

payload := managedenvironmentsstorages.ManagedEnvironmentStorage{
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


### Example Usage: `ManagedEnvironmentsStoragesClient.Delete`

```go
ctx := context.TODO()
id := managedenvironmentsstorages.NewStorageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName", "storageName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedEnvironmentsStoragesClient.Get`

```go
ctx := context.TODO()
id := managedenvironmentsstorages.NewStorageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName", "storageName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedEnvironmentsStoragesClient.List`

```go
ctx := context.TODO()
id := managedenvironmentsstorages.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
