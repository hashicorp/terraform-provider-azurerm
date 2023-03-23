
## `github.com/hashicorp/go-azure-sdk/resource-manager/deviceupdate/2022-10-01/deviceupdates` Documentation

The `deviceupdates` SDK allows for interaction with the Azure Resource Manager Service `deviceupdate` (API Version `2022-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/deviceupdate/2022-10-01/deviceupdates"
```


### Client Initialization

```go
client := deviceupdates.NewDeviceupdatesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeviceupdatesClient.AccountsCreate`

```go
ctx := context.TODO()
id := deviceupdates.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := deviceupdates.Account{
	// ...
}


if err := client.AccountsCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeviceupdatesClient.AccountsDelete`

```go
ctx := context.TODO()
id := deviceupdates.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

if err := client.AccountsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DeviceupdatesClient.AccountsGet`

```go
ctx := context.TODO()
id := deviceupdates.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.AccountsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeviceupdatesClient.AccountsHead`

```go
ctx := context.TODO()
id := deviceupdates.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.AccountsHead(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeviceupdatesClient.AccountsListByResourceGroup`

```go
ctx := context.TODO()
id := deviceupdates.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.AccountsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.AccountsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeviceupdatesClient.AccountsListBySubscription`

```go
ctx := context.TODO()
id := deviceupdates.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.AccountsListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.AccountsListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeviceupdatesClient.AccountsUpdate`

```go
ctx := context.TODO()
id := deviceupdates.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := deviceupdates.AccountUpdate{
	// ...
}


if err := client.AccountsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeviceupdatesClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := deviceupdates.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := deviceupdates.CheckNameAvailabilityRequest{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeviceupdatesClient.InstancesCreate`

```go
ctx := context.TODO()
id := deviceupdates.NewInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "instanceValue")

payload := deviceupdates.Instance{
	// ...
}


if err := client.InstancesCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DeviceupdatesClient.InstancesDelete`

```go
ctx := context.TODO()
id := deviceupdates.NewInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "instanceValue")

if err := client.InstancesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DeviceupdatesClient.InstancesGet`

```go
ctx := context.TODO()
id := deviceupdates.NewInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "instanceValue")

read, err := client.InstancesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeviceupdatesClient.InstancesHead`

```go
ctx := context.TODO()
id := deviceupdates.NewInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "instanceValue")

read, err := client.InstancesHead(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeviceupdatesClient.InstancesListByAccount`

```go
ctx := context.TODO()
id := deviceupdates.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

// alternatively `client.InstancesListByAccount(ctx, id)` can be used to do batched pagination
items, err := client.InstancesListByAccountComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeviceupdatesClient.InstancesUpdate`

```go
ctx := context.TODO()
id := deviceupdates.NewInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "instanceValue")

payload := deviceupdates.TagUpdate{
	// ...
}


read, err := client.InstancesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeviceupdatesClient.PrivateEndpointConnectionProxiesUpdatePrivateEndpointProperties`

```go
ctx := context.TODO()
id := deviceupdates.NewPrivateEndpointConnectionProxyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "privateEndpointConnectionProxyIdValue")

payload := deviceupdates.PrivateEndpointUpdate{
	// ...
}


read, err := client.PrivateEndpointConnectionProxiesUpdatePrivateEndpointProperties(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeviceupdatesClient.PrivateEndpointConnectionProxiesValidate`

```go
ctx := context.TODO()
id := deviceupdates.NewPrivateEndpointConnectionProxyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "privateEndpointConnectionProxyIdValue")

payload := deviceupdates.PrivateEndpointConnectionProxy{
	// ...
}


read, err := client.PrivateEndpointConnectionProxiesValidate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
