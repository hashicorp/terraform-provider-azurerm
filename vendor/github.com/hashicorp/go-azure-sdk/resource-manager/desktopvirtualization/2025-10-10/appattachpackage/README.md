
## `github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/appattachpackage` Documentation

The `appattachpackage` SDK allows for interaction with Azure Resource Manager `desktopvirtualization` (API Version `2025-10-10`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/appattachpackage"
```


### Client Initialization

```go
client := appattachpackage.NewAppAttachPackageClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AppAttachPackageClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := appattachpackage.NewAppAttachPackageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "appAttachPackageName")

payload := appattachpackage.AppAttachPackage{
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


### Example Usage: `AppAttachPackageClient.Delete`

```go
ctx := context.TODO()
id := appattachpackage.NewAppAttachPackageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "appAttachPackageName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppAttachPackageClient.Get`

```go
ctx := context.TODO()
id := appattachpackage.NewAppAttachPackageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "appAttachPackageName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppAttachPackageClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, appattachpackage.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, appattachpackage.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppAttachPackageClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, appattachpackage.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, appattachpackage.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppAttachPackageClient.Update`

```go
ctx := context.TODO()
id := appattachpackage.NewAppAttachPackageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "appAttachPackageName")

payload := appattachpackage.AppAttachPackagePatch{
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
