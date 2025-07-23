
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/attachednetworkconnections` Documentation

The `attachednetworkconnections` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/attachednetworkconnections"
```


### Client Initialization

```go
client := attachednetworkconnections.NewAttachedNetworkConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AttachedNetworkConnectionsClient.AttachedNetworksCreateOrUpdate`

```go
ctx := context.TODO()
id := attachednetworkconnections.NewDevCenterAttachedNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "attachedNetworkName")

payload := attachednetworkconnections.AttachedNetworkConnection{
	// ...
}


if err := client.AttachedNetworksCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AttachedNetworkConnectionsClient.AttachedNetworksDelete`

```go
ctx := context.TODO()
id := attachednetworkconnections.NewDevCenterAttachedNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "attachedNetworkName")

if err := client.AttachedNetworksDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AttachedNetworkConnectionsClient.AttachedNetworksGetByDevCenter`

```go
ctx := context.TODO()
id := attachednetworkconnections.NewDevCenterAttachedNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "attachedNetworkName")

read, err := client.AttachedNetworksGetByDevCenter(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AttachedNetworkConnectionsClient.AttachedNetworksGetByProject`

```go
ctx := context.TODO()
id := attachednetworkconnections.NewAttachedNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "attachedNetworkName")

read, err := client.AttachedNetworksGetByProject(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AttachedNetworkConnectionsClient.AttachedNetworksListByDevCenter`

```go
ctx := context.TODO()
id := attachednetworkconnections.NewDevCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName")

// alternatively `client.AttachedNetworksListByDevCenter(ctx, id, attachednetworkconnections.DefaultAttachedNetworksListByDevCenterOperationOptions())` can be used to do batched pagination
items, err := client.AttachedNetworksListByDevCenterComplete(ctx, id, attachednetworkconnections.DefaultAttachedNetworksListByDevCenterOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AttachedNetworkConnectionsClient.AttachedNetworksListByProject`

```go
ctx := context.TODO()
id := attachednetworkconnections.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName")

// alternatively `client.AttachedNetworksListByProject(ctx, id, attachednetworkconnections.DefaultAttachedNetworksListByProjectOperationOptions())` can be used to do batched pagination
items, err := client.AttachedNetworksListByProjectComplete(ctx, id, attachednetworkconnections.DefaultAttachedNetworksListByProjectOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
