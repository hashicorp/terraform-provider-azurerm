
## `github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-05-01/volumes` Documentation

The `volumes` SDK allows for interaction with the Azure Resource Manager Service `netapp` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-05-01/volumes"
```


### Client Initialization

```go
client := volumes.NewVolumesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VolumesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "capacityPoolValue", "volumeValue")

payload := volumes.Volume{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.Delete`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "capacityPoolValue", "volumeValue")

if err := client.DeleteThenPoll(ctx, id, volumes.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.Get`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "capacityPoolValue", "volumeValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VolumesClient.List`

```go
ctx := context.TODO()
id := volumes.NewCapacityPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "capacityPoolValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VolumesClient.Update`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "capacityPoolValue", "volumeValue")

payload := volumes.VolumePatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
