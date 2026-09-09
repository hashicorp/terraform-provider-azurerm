
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/privatelinkservices` Documentation

The `privatelinkservices` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/privatelinkservices"
```


### Client Initialization

```go
client := privatelinkservices.NewPrivateLinkServicesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateLinkServicesClient.CheckPrivateLinkServiceVisibility`

```go
ctx := context.TODO()
id := privatelinkservices.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := privatelinkservices.CheckPrivateLinkServiceVisibilityRequest{
	// ...
}


if err := client.CheckPrivateLinkServiceVisibilityThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateLinkServicesClient.CheckPrivateLinkServiceVisibilityByResourceGroup`

```go
ctx := context.TODO()
id := privatelinkservices.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName")

payload := privatelinkservices.CheckPrivateLinkServiceVisibilityRequest{
	// ...
}


if err := client.CheckPrivateLinkServiceVisibilityByResourceGroupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateLinkServicesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := privatelinkservices.NewPrivateLinkServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkServiceName")

payload := privatelinkservices.PrivateLinkService{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateLinkServicesClient.Delete`

```go
ctx := context.TODO()
id := privatelinkservices.NewPrivateLinkServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkServiceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateLinkServicesClient.DeletePrivateEndpointConnection`

```go
ctx := context.TODO()
id := privatelinkservices.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkServiceName", "privateEndpointConnectionName")

if err := client.DeletePrivateEndpointConnectionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateLinkServicesClient.Get`

```go
ctx := context.TODO()
id := privatelinkservices.NewPrivateLinkServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkServiceName")

read, err := client.Get(ctx, id, privatelinkservices.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkServicesClient.GetPrivateEndpointConnection`

```go
ctx := context.TODO()
id := privatelinkservices.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkServiceName", "privateEndpointConnectionName")

read, err := client.GetPrivateEndpointConnection(ctx, id, privatelinkservices.DefaultGetPrivateEndpointConnectionOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkServicesClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateLinkServicesClient.ListAutoApprovedPrivateLinkServices`

```go
ctx := context.TODO()
id := privatelinkservices.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.ListAutoApprovedPrivateLinkServices(ctx, id)` can be used to do batched pagination
items, err := client.ListAutoApprovedPrivateLinkServicesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateLinkServicesClient.ListAutoApprovedPrivateLinkServicesByResourceGroup`

```go
ctx := context.TODO()
id := privatelinkservices.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName")

// alternatively `client.ListAutoApprovedPrivateLinkServicesByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListAutoApprovedPrivateLinkServicesByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateLinkServicesClient.ListBySubscription`

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


### Example Usage: `PrivateLinkServicesClient.ListPrivateEndpointConnections`

```go
ctx := context.TODO()
id := privatelinkservices.NewPrivateLinkServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkServiceName")

// alternatively `client.ListPrivateEndpointConnections(ctx, id)` can be used to do batched pagination
items, err := client.ListPrivateEndpointConnectionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateLinkServicesClient.UpdatePrivateEndpointConnection`

```go
ctx := context.TODO()
id := privatelinkservices.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkServiceName", "privateEndpointConnectionName")

payload := privatelinkservices.PrivateEndpointConnection{
	// ...
}


read, err := client.UpdatePrivateEndpointConnection(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
