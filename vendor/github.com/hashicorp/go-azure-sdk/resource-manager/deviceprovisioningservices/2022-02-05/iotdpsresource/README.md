
## `github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource` Documentation

The `iotdpsresource` SDK allows for interaction with Azure Resource Manager `deviceprovisioningservices` (API Version `2022-02-05`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
```


### Client Initialization

```go
client := iotdpsresource.NewIotDpsResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IotDpsResourceClient.CheckProvisioningServiceNameAvailability`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := iotdpsresource.OperationInputs{
	// ...
}


read, err := client.CheckProvisioningServiceNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IotDpsResourceClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewProvisioningServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName")

payload := iotdpsresource.ProvisioningServiceDescription{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `IotDpsResourceClient.CreateOrUpdatePrivateEndpointConnection`

```go
ctx := context.TODO()
id := iotdpsresource.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName", "privateEndpointConnectionName")

payload := iotdpsresource.PrivateEndpointConnection{
	// ...
}


if err := client.CreateOrUpdatePrivateEndpointConnectionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `IotDpsResourceClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewProvisioningServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IotDpsResourceClient.DeletePrivateEndpointConnection`

```go
ctx := context.TODO()
id := iotdpsresource.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName", "privateEndpointConnectionName")

if err := client.DeletePrivateEndpointConnectionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IotDpsResourceClient.Get`

```go
ctx := context.TODO()
id := commonids.NewProvisioningServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IotDpsResourceClient.GetPrivateEndpointConnection`

```go
ctx := context.TODO()
id := iotdpsresource.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName", "privateEndpointConnectionName")

read, err := client.GetPrivateEndpointConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IotDpsResourceClient.GetPrivateLinkResources`

```go
ctx := context.TODO()
id := iotdpsresource.NewPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName", "groupId")

read, err := client.GetPrivateLinkResources(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IotDpsResourceClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IotDpsResourceClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IotDpsResourceClient.ListKeys`

```go
ctx := context.TODO()
id := commonids.NewProvisioningServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName")

// alternatively `client.ListKeys(ctx, id)` can be used to do batched pagination
items, err := client.ListKeysComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IotDpsResourceClient.ListKeysForKeyName`

```go
ctx := context.TODO()
id := iotdpsresource.NewKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName", "keyName")

read, err := client.ListKeysForKeyName(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IotDpsResourceClient.ListPrivateEndpointConnections`

```go
ctx := context.TODO()
id := commonids.NewProvisioningServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName")

read, err := client.ListPrivateEndpointConnections(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IotDpsResourceClient.ListPrivateLinkResources`

```go
ctx := context.TODO()
id := commonids.NewProvisioningServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName")

read, err := client.ListPrivateLinkResources(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IotDpsResourceClient.ListValidSkus`

```go
ctx := context.TODO()
id := commonids.NewProvisioningServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName")

// alternatively `client.ListValidSkus(ctx, id)` can be used to do batched pagination
items, err := client.ListValidSkusComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IotDpsResourceClient.Update`

```go
ctx := context.TODO()
id := commonids.NewProvisioningServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName")

payload := iotdpsresource.TagsResource{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
