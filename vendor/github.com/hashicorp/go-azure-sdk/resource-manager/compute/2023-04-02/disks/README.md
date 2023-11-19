
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks` Documentation

The `disks` SDK allows for interaction with the Azure Resource Manager Service `compute` (API Version `2023-04-02`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
```


### Client Initialization

```go
client := disks.NewDisksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DisksClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := disks.NewDiskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskValue")

payload := disks.Disk{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DisksClient.Delete`

```go
ctx := context.TODO()
id := disks.NewDiskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DisksClient.Get`

```go
ctx := context.TODO()
id := disks.NewDiskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DisksClient.GrantAccess`

```go
ctx := context.TODO()
id := disks.NewDiskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskValue")

payload := disks.GrantAccessData{
	// ...
}


if err := client.GrantAccessThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DisksClient.List`

```go
ctx := context.TODO()
id := disks.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DisksClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := disks.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DisksClient.RevokeAccess`

```go
ctx := context.TODO()
id := disks.NewDiskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskValue")

if err := client.RevokeAccessThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DisksClient.Update`

```go
ctx := context.TODO()
id := disks.NewDiskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskValue")

payload := disks.DiskUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
