
## `github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/accounts` Documentation

The `accounts` SDK allows for interaction with the Azure Resource Manager Service `media` (API Version `2020-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/accounts"
```


### Client Initialization

```go
client := accounts.NewAccountsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AccountsClient.LocationsCheckNameAvailability`

```go
ctx := context.TODO()
id := accounts.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

payload := accounts.CheckNameAvailabilityInput{
	// ...
}


read, err := client.LocationsCheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.MediaservicesCreateOrUpdate`

```go
ctx := context.TODO()
id := accounts.NewProviderMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

payload := accounts.MediaService{
	// ...
}


read, err := client.MediaservicesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.MediaservicesDelete`

```go
ctx := context.TODO()
id := accounts.NewProviderMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

read, err := client.MediaservicesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.MediaservicesGet`

```go
ctx := context.TODO()
id := accounts.NewProviderMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

read, err := client.MediaservicesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.MediaservicesGetBySubscription`

```go
ctx := context.TODO()
id := accounts.NewMediaServiceID("12345678-1234-9876-4563-123456789012", "mediaServiceValue")

read, err := client.MediaservicesGetBySubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.MediaservicesList`

```go
ctx := context.TODO()
id := accounts.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.MediaservicesList(ctx, id)` can be used to do batched pagination
items, err := client.MediaservicesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AccountsClient.MediaservicesListBySubscription`

```go
ctx := context.TODO()
id := accounts.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.MediaservicesListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.MediaservicesListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AccountsClient.MediaservicesListEdgePolicies`

```go
ctx := context.TODO()
id := accounts.NewProviderMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

payload := accounts.ListEdgePoliciesInput{
	// ...
}


read, err := client.MediaservicesListEdgePolicies(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.MediaservicesSyncStorageKeys`

```go
ctx := context.TODO()
id := accounts.NewProviderMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

payload := accounts.SyncStorageKeysInput{
	// ...
}


read, err := client.MediaservicesSyncStorageKeys(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.MediaservicesUpdate`

```go
ctx := context.TODO()
id := accounts.NewProviderMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

payload := accounts.MediaService{
	// ...
}


read, err := client.MediaservicesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.PrivateEndpointConnectionsCreateOrUpdate`

```go
ctx := context.TODO()
id := accounts.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "privateEndpointConnectionValue")

payload := accounts.PrivateEndpointConnection{
	// ...
}


read, err := client.PrivateEndpointConnectionsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.PrivateEndpointConnectionsDelete`

```go
ctx := context.TODO()
id := accounts.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "privateEndpointConnectionValue")

read, err := client.PrivateEndpointConnectionsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.PrivateEndpointConnectionsGet`

```go
ctx := context.TODO()
id := accounts.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "privateEndpointConnectionValue")

read, err := client.PrivateEndpointConnectionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.PrivateEndpointConnectionsList`

```go
ctx := context.TODO()
id := accounts.NewProviderMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

read, err := client.PrivateEndpointConnectionsList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.PrivateLinkResourcesGet`

```go
ctx := context.TODO()
id := accounts.NewPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "privateLinkResourceValue")

read, err := client.PrivateLinkResourcesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.PrivateLinkResourcesList`

```go
ctx := context.TODO()
id := accounts.NewProviderMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

read, err := client.PrivateLinkResourcesList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
