
## `github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/msixpackage` Documentation

The `msixpackage` SDK allows for interaction with Azure Resource Manager `desktopvirtualization` (API Version `2022-02-10-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/msixpackage"
```


### Client Initialization

```go
client := msixpackage.NewMSIXPackageClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MSIXPackageClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := msixpackage.NewMsixPackageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolName", "msixPackageName")

payload := msixpackage.MSIXPackage{
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


### Example Usage: `MSIXPackageClient.Delete`

```go
ctx := context.TODO()
id := msixpackage.NewMsixPackageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolName", "msixPackageName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MSIXPackageClient.Get`

```go
ctx := context.TODO()
id := msixpackage.NewMsixPackageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolName", "msixPackageName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MSIXPackageClient.List`

```go
ctx := context.TODO()
id := msixpackage.NewHostPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MSIXPackageClient.Update`

```go
ctx := context.TODO()
id := msixpackage.NewMsixPackageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolName", "msixPackageName")

payload := msixpackage.MSIXPackagePatch{
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
