
## `github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-05-01/volumegroups` Documentation

The `volumegroups` SDK allows for interaction with the Azure Resource Manager Service `netapp` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-05-01/volumegroups"
```


### Client Initialization

```go
client := volumegroups.NewVolumeGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VolumeGroupsClient.Create`

```go
ctx := context.TODO()
id := volumegroups.NewVolumeGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "volumeGroupValue")

payload := volumegroups.VolumeGroupDetails{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumeGroupsClient.Delete`

```go
ctx := context.TODO()
id := volumegroups.NewVolumeGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "volumeGroupValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumeGroupsClient.Get`

```go
ctx := context.TODO()
id := volumegroups.NewVolumeGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "volumeGroupValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VolumeGroupsClient.ListByNetAppAccount`

```go
ctx := context.TODO()
id := volumegroups.NewNetAppAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue")

read, err := client.ListByNetAppAccount(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
