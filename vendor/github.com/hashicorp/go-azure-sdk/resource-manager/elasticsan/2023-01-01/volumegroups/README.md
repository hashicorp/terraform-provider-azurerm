
## `github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumegroups` Documentation

The `volumegroups` SDK allows for interaction with Azure Resource Manager `elasticsan` (API Version `2023-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumegroups"
```


### Client Initialization

```go
client := volumegroups.NewVolumeGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VolumeGroupsClient.Create`

```go
ctx := context.TODO()
id := volumegroups.NewVolumeGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "elasticSanName", "volumeGroupName")

payload := volumegroups.VolumeGroup{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumeGroupsClient.Delete`

```go
ctx := context.TODO()
id := volumegroups.NewVolumeGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "elasticSanName", "volumeGroupName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumeGroupsClient.Get`

```go
ctx := context.TODO()
id := volumegroups.NewVolumeGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "elasticSanName", "volumeGroupName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VolumeGroupsClient.ListByElasticSan`

```go
ctx := context.TODO()
id := volumegroups.NewElasticSanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "elasticSanName")

// alternatively `client.ListByElasticSan(ctx, id)` can be used to do batched pagination
items, err := client.ListByElasticSanComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VolumeGroupsClient.Update`

```go
ctx := context.TODO()
id := volumegroups.NewVolumeGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "elasticSanName", "volumeGroupName")

payload := volumegroups.VolumeGroupUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
