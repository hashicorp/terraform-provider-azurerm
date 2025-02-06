
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/storagesyncservicesresource` Documentation

The `storagesyncservicesresource` SDK allows for interaction with Azure Resource Manager `storagesync` (API Version `2020-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/storagesyncservicesresource"
```


### Client Initialization

```go
client := storagesyncservicesresource.NewStorageSyncServicesResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StorageSyncServicesResourceClient.StorageSyncServicesCreate`

```go
ctx := context.TODO()
id := storagesyncservicesresource.NewStorageSyncServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName")

payload := storagesyncservicesresource.StorageSyncServiceCreateParameters{
	// ...
}


if err := client.StorageSyncServicesCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StorageSyncServicesResourceClient.StorageSyncServicesDelete`

```go
ctx := context.TODO()
id := storagesyncservicesresource.NewStorageSyncServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName")

if err := client.StorageSyncServicesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StorageSyncServicesResourceClient.StorageSyncServicesGet`

```go
ctx := context.TODO()
id := storagesyncservicesresource.NewStorageSyncServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName")

read, err := client.StorageSyncServicesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageSyncServicesResourceClient.StorageSyncServicesListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.StorageSyncServicesListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageSyncServicesResourceClient.StorageSyncServicesListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.StorageSyncServicesListBySubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageSyncServicesResourceClient.StorageSyncServicesUpdate`

```go
ctx := context.TODO()
id := storagesyncservicesresource.NewStorageSyncServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName")

payload := storagesyncservicesresource.StorageSyncServiceUpdateParameters{
	// ...
}


if err := client.StorageSyncServicesUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
