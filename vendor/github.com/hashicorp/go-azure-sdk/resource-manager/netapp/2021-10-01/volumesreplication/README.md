
## `github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/volumesreplication` Documentation

The `volumesreplication` SDK allows for interaction with the Azure Resource Manager Service `netapp` (API Version `2021-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/volumesreplication"
```


### Client Initialization

```go
client := volumesreplication.NewVolumesReplicationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VolumesReplicationClient.VolumesAuthorizeReplication`

```go
ctx := context.TODO()
id := volumesreplication.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "poolValue", "volumeValue")

payload := volumesreplication.AuthorizeRequest{
	// ...
}


if err := client.VolumesAuthorizeReplicationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesReplicationClient.VolumesBreakReplication`

```go
ctx := context.TODO()
id := volumesreplication.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "poolValue", "volumeValue")

payload := volumesreplication.BreakReplicationRequest{
	// ...
}


if err := client.VolumesBreakReplicationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesReplicationClient.VolumesDeleteReplication`

```go
ctx := context.TODO()
id := volumesreplication.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "poolValue", "volumeValue")

if err := client.VolumesDeleteReplicationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesReplicationClient.VolumesReInitializeReplication`

```go
ctx := context.TODO()
id := volumesreplication.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "poolValue", "volumeValue")

if err := client.VolumesReInitializeReplicationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesReplicationClient.VolumesReplicationStatus`

```go
ctx := context.TODO()
id := volumesreplication.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "poolValue", "volumeValue")

read, err := client.VolumesReplicationStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VolumesReplicationClient.VolumesResyncReplication`

```go
ctx := context.TODO()
id := volumesreplication.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "poolValue", "volumeValue")

if err := client.VolumesResyncReplicationThenPoll(ctx, id); err != nil {
	// handle the error
}
```
