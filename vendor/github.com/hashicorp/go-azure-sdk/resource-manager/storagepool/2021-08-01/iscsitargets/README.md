
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/iscsitargets` Documentation

The `iscsitargets` SDK allows for interaction with the Azure Resource Manager Service `storagepool` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/iscsitargets"
```


### Client Initialization

```go
client := iscsitargets.NewIscsiTargetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IscsiTargetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := iscsitargets.NewIscsiTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue", "iscsiTargetValue")

payload := iscsitargets.IscsiTargetCreate{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `IscsiTargetsClient.Delete`

```go
ctx := context.TODO()
id := iscsitargets.NewIscsiTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue", "iscsiTargetValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IscsiTargetsClient.Get`

```go
ctx := context.TODO()
id := iscsitargets.NewIscsiTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue", "iscsiTargetValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IscsiTargetsClient.ListByDiskPool`

```go
ctx := context.TODO()
id := iscsitargets.NewDiskPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue")

// alternatively `client.ListByDiskPool(ctx, id)` can be used to do batched pagination
items, err := client.ListByDiskPoolComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IscsiTargetsClient.Update`

```go
ctx := context.TODO()
id := iscsitargets.NewIscsiTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue", "iscsiTargetValue")

payload := iscsitargets.IscsiTargetUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
