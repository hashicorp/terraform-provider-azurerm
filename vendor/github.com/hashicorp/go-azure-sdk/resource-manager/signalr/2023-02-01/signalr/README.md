
## `github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2023-02-01/signalr` Documentation

The `signalr` SDK allows for interaction with the Azure Resource Manager Service `signalr` (API Version `2023-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2023-02-01/signalr"
```


### Client Initialization

```go
client := signalr.NewSignalRClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SignalRClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := signalr.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

payload := signalr.NameAvailabilityParameters{
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


### Example Usage: `SignalRClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := signalr.NewSignalRID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue")

payload := signalr.SignalRResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SignalRClient.CustomCertificatesCreateOrUpdate`

```go
ctx := context.TODO()
id := signalr.NewCustomCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue", "customCertificateValue")

payload := signalr.CustomCertificate{
	// ...
}


if err := client.CustomCertificatesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SignalRClient.CustomCertificatesDelete`

```go
ctx := context.TODO()
id := signalr.NewCustomCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue", "customCertificateValue")

read, err := client.CustomCertificatesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SignalRClient.CustomCertificatesGet`

```go
ctx := context.TODO()
id := signalr.NewCustomCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue", "customCertificateValue")

read, err := client.CustomCertificatesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SignalRClient.CustomCertificatesList`

```go
ctx := context.TODO()
id := signalr.NewSignalRID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue")

// alternatively `client.CustomCertificatesList(ctx, id)` can be used to do batched pagination
items, err := client.CustomCertificatesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SignalRClient.CustomDomainsCreateOrUpdate`

```go
ctx := context.TODO()
id := signalr.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue", "customDomainValue")

payload := signalr.CustomDomain{
	// ...
}


if err := client.CustomDomainsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SignalRClient.CustomDomainsDelete`

```go
ctx := context.TODO()
id := signalr.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue", "customDomainValue")

if err := client.CustomDomainsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SignalRClient.CustomDomainsGet`

```go
ctx := context.TODO()
id := signalr.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue", "customDomainValue")

read, err := client.CustomDomainsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SignalRClient.CustomDomainsList`

```go
ctx := context.TODO()
id := signalr.NewSignalRID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue")

// alternatively `client.CustomDomainsList(ctx, id)` can be used to do batched pagination
items, err := client.CustomDomainsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SignalRClient.Delete`

```go
ctx := context.TODO()
id := signalr.NewSignalRID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SignalRClient.Get`

```go
ctx := context.TODO()
id := signalr.NewSignalRID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SignalRClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := signalr.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SignalRClient.ListBySubscription`

```go
ctx := context.TODO()
id := signalr.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SignalRClient.ListKeys`

```go
ctx := context.TODO()
id := signalr.NewSignalRID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue")

read, err := client.ListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SignalRClient.ListSkus`

```go
ctx := context.TODO()
id := signalr.NewSignalRID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue")

read, err := client.ListSkus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SignalRClient.PrivateEndpointConnectionsDelete`

```go
ctx := context.TODO()
id := signalr.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue", "privateEndpointConnectionValue")

if err := client.PrivateEndpointConnectionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SignalRClient.PrivateEndpointConnectionsGet`

```go
ctx := context.TODO()
id := signalr.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue", "privateEndpointConnectionValue")

read, err := client.PrivateEndpointConnectionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SignalRClient.PrivateEndpointConnectionsList`

```go
ctx := context.TODO()
id := signalr.NewSignalRID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue")

// alternatively `client.PrivateEndpointConnectionsList(ctx, id)` can be used to do batched pagination
items, err := client.PrivateEndpointConnectionsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SignalRClient.PrivateEndpointConnectionsUpdate`

```go
ctx := context.TODO()
id := signalr.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue", "privateEndpointConnectionValue")

payload := signalr.PrivateEndpointConnection{
	// ...
}


read, err := client.PrivateEndpointConnectionsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SignalRClient.PrivateLinkResourcesList`

```go
ctx := context.TODO()
id := signalr.NewSignalRID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue")

// alternatively `client.PrivateLinkResourcesList(ctx, id)` can be used to do batched pagination
items, err := client.PrivateLinkResourcesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SignalRClient.RegenerateKey`

```go
ctx := context.TODO()
id := signalr.NewSignalRID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue")

payload := signalr.RegenerateKeyParameters{
	// ...
}


if err := client.RegenerateKeyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SignalRClient.Restart`

```go
ctx := context.TODO()
id := signalr.NewSignalRID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue")

if err := client.RestartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SignalRClient.SharedPrivateLinkResourcesCreateOrUpdate`

```go
ctx := context.TODO()
id := signalr.NewSharedPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue", "sharedPrivateLinkResourceValue")

payload := signalr.SharedPrivateLinkResource{
	// ...
}


if err := client.SharedPrivateLinkResourcesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SignalRClient.SharedPrivateLinkResourcesDelete`

```go
ctx := context.TODO()
id := signalr.NewSharedPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue", "sharedPrivateLinkResourceValue")

if err := client.SharedPrivateLinkResourcesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SignalRClient.SharedPrivateLinkResourcesGet`

```go
ctx := context.TODO()
id := signalr.NewSharedPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue", "sharedPrivateLinkResourceValue")

read, err := client.SharedPrivateLinkResourcesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SignalRClient.SharedPrivateLinkResourcesList`

```go
ctx := context.TODO()
id := signalr.NewSignalRID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue")

// alternatively `client.SharedPrivateLinkResourcesList(ctx, id)` can be used to do batched pagination
items, err := client.SharedPrivateLinkResourcesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SignalRClient.Update`

```go
ctx := context.TODO()
id := signalr.NewSignalRID("12345678-1234-9876-4563-123456789012", "example-resource-group", "signalRValue")

payload := signalr.SignalRResource{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SignalRClient.UsagesList`

```go
ctx := context.TODO()
id := signalr.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

// alternatively `client.UsagesList(ctx, id)` can be used to do batched pagination
items, err := client.UsagesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
