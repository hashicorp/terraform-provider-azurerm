
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkmanagerconnections` Documentation

The `networkmanagerconnections` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkmanagerconnections"
```


### Client Initialization

```go
client := networkmanagerconnections.NewNetworkManagerConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkManagerConnectionsClient.ManagementGroupNetworkManagerConnectionsCreateOrUpdate`

```go
ctx := context.TODO()
id := networkmanagerconnections.NewProviders2NetworkManagerConnectionID("managementGroupId", "networkManagerConnectionName")

payload := networkmanagerconnections.NetworkManagerConnection{
	// ...
}


read, err := client.ManagementGroupNetworkManagerConnectionsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkManagerConnectionsClient.ManagementGroupNetworkManagerConnectionsDelete`

```go
ctx := context.TODO()
id := networkmanagerconnections.NewProviders2NetworkManagerConnectionID("managementGroupId", "networkManagerConnectionName")

read, err := client.ManagementGroupNetworkManagerConnectionsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkManagerConnectionsClient.ManagementGroupNetworkManagerConnectionsGet`

```go
ctx := context.TODO()
id := networkmanagerconnections.NewProviders2NetworkManagerConnectionID("managementGroupId", "networkManagerConnectionName")

read, err := client.ManagementGroupNetworkManagerConnectionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkManagerConnectionsClient.ManagementGroupNetworkManagerConnectionsList`

```go
ctx := context.TODO()
id := commonids.NewManagementGroupID("groupId")

// alternatively `client.ManagementGroupNetworkManagerConnectionsList(ctx, id, networkmanagerconnections.DefaultManagementGroupNetworkManagerConnectionsListOperationOptions())` can be used to do batched pagination
items, err := client.ManagementGroupNetworkManagerConnectionsListComplete(ctx, id, networkmanagerconnections.DefaultManagementGroupNetworkManagerConnectionsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkManagerConnectionsClient.SubscriptionNetworkManagerConnectionsCreateOrUpdate`

```go
ctx := context.TODO()
id := networkmanagerconnections.NewNetworkManagerConnectionID("12345678-1234-9876-4563-123456789012", "networkManagerConnectionName")

payload := networkmanagerconnections.NetworkManagerConnection{
	// ...
}


read, err := client.SubscriptionNetworkManagerConnectionsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkManagerConnectionsClient.SubscriptionNetworkManagerConnectionsDelete`

```go
ctx := context.TODO()
id := networkmanagerconnections.NewNetworkManagerConnectionID("12345678-1234-9876-4563-123456789012", "networkManagerConnectionName")

read, err := client.SubscriptionNetworkManagerConnectionsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkManagerConnectionsClient.SubscriptionNetworkManagerConnectionsGet`

```go
ctx := context.TODO()
id := networkmanagerconnections.NewNetworkManagerConnectionID("12345678-1234-9876-4563-123456789012", "networkManagerConnectionName")

read, err := client.SubscriptionNetworkManagerConnectionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkManagerConnectionsClient.SubscriptionNetworkManagerConnectionsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.SubscriptionNetworkManagerConnectionsList(ctx, id, networkmanagerconnections.DefaultSubscriptionNetworkManagerConnectionsListOperationOptions())` can be used to do batched pagination
items, err := client.SubscriptionNetworkManagerConnectionsListComplete(ctx, id, networkmanagerconnections.DefaultSubscriptionNetworkManagerConnectionsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
