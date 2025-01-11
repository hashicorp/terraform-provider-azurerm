
## `github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/attacheddatanetwork` Documentation

The `attacheddatanetwork` SDK allows for interaction with Azure Resource Manager `mobilenetwork` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/attacheddatanetwork"
```


### Client Initialization

```go
client := attacheddatanetwork.NewAttachedDataNetworkClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AttachedDataNetworkClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := attacheddatanetwork.NewAttachedDataNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "packetCoreControlPlaneName", "packetCoreDataPlaneName", "attachedDataNetworkName")

payload := attacheddatanetwork.AttachedDataNetwork{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AttachedDataNetworkClient.Delete`

```go
ctx := context.TODO()
id := attacheddatanetwork.NewAttachedDataNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "packetCoreControlPlaneName", "packetCoreDataPlaneName", "attachedDataNetworkName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AttachedDataNetworkClient.Get`

```go
ctx := context.TODO()
id := attacheddatanetwork.NewAttachedDataNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "packetCoreControlPlaneName", "packetCoreDataPlaneName", "attachedDataNetworkName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AttachedDataNetworkClient.UpdateTags`

```go
ctx := context.TODO()
id := attacheddatanetwork.NewAttachedDataNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "packetCoreControlPlaneName", "packetCoreDataPlaneName", "attachedDataNetworkName")

payload := attacheddatanetwork.TagsObject{
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
