
## `github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcorecontrolplane` Documentation

The `packetcorecontrolplane` SDK allows for interaction with Azure Resource Manager `mobilenetwork` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcorecontrolplane"
```


### Client Initialization

```go
client := packetcorecontrolplane.NewPacketCoreControlPlaneClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PacketCoreControlPlaneClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := packetcorecontrolplane.NewPacketCoreControlPlaneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "packetCoreControlPlaneName")

payload := packetcorecontrolplane.PacketCoreControlPlane{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PacketCoreControlPlaneClient.Delete`

```go
ctx := context.TODO()
id := packetcorecontrolplane.NewPacketCoreControlPlaneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "packetCoreControlPlaneName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PacketCoreControlPlaneClient.Get`

```go
ctx := context.TODO()
id := packetcorecontrolplane.NewPacketCoreControlPlaneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "packetCoreControlPlaneName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PacketCoreControlPlaneClient.UpdateTags`

```go
ctx := context.TODO()
id := packetcorecontrolplane.NewPacketCoreControlPlaneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "packetCoreControlPlaneName")

payload := packetcorecontrolplane.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
