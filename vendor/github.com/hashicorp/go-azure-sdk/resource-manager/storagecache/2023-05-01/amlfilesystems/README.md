
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/amlfilesystems` Documentation

The `amlfilesystems` SDK allows for interaction with Azure Resource Manager `storagecache` (API Version `2023-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/amlfilesystems"
```


### Client Initialization

```go
client := amlfilesystems.NewAmlFilesystemsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AmlFilesystemsClient.Archive`

```go
ctx := context.TODO()
id := amlfilesystems.NewAmlFilesystemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName")

payload := amlfilesystems.AmlFilesystemArchiveInfo{
	// ...
}


read, err := client.Archive(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AmlFilesystemsClient.CancelArchive`

```go
ctx := context.TODO()
id := amlfilesystems.NewAmlFilesystemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName")

read, err := client.CancelArchive(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AmlFilesystemsClient.CheckAmlFSSubnets`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := amlfilesystems.AmlFilesystemSubnetInfo{
	// ...
}


read, err := client.CheckAmlFSSubnets(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AmlFilesystemsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := amlfilesystems.NewAmlFilesystemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName")

payload := amlfilesystems.AmlFilesystem{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AmlFilesystemsClient.Delete`

```go
ctx := context.TODO()
id := amlfilesystems.NewAmlFilesystemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AmlFilesystemsClient.Get`

```go
ctx := context.TODO()
id := amlfilesystems.NewAmlFilesystemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AmlFilesystemsClient.GetRequiredAmlFSSubnetsSize`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := amlfilesystems.RequiredAmlFilesystemSubnetsSizeInfo{
	// ...
}


read, err := client.GetRequiredAmlFSSubnetsSize(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AmlFilesystemsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AmlFilesystemsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AmlFilesystemsClient.Update`

```go
ctx := context.TODO()
id := amlfilesystems.NewAmlFilesystemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName")

payload := amlfilesystems.AmlFilesystemUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
