
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/attachednetworkconnections` Documentation

The `attachednetworkconnections` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/attachednetworkconnections"
```


### Client Initialization

```go
client := attachednetworkconnections.NewAttachedNetworkConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
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
