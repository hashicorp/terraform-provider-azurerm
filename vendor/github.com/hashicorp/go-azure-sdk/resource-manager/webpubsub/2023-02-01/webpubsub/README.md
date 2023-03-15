
## `github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2023-02-01/webpubsub` Documentation

The `webpubsub` SDK allows for interaction with the Azure Resource Manager Service `webpubsub` (API Version `2023-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2023-02-01/webpubsub"
```


### Client Initialization

```go
client := webpubsub.NewWebPubSubClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `WebPubSubClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := webpubsub.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

payload := webpubsub.NameAvailabilityParameters{
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


### Example Usage: `WebPubSubClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue")

payload := webpubsub.WebPubSubResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.CustomCertificatesCreateOrUpdate`

```go
ctx := context.TODO()
id := webpubsub.NewCustomCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue", "customCertificateValue")

payload := webpubsub.CustomCertificate{
	// ...
}


if err := client.CustomCertificatesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.CustomCertificatesDelete`

```go
ctx := context.TODO()
id := webpubsub.NewCustomCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue", "customCertificateValue")

read, err := client.CustomCertificatesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebPubSubClient.CustomCertificatesGet`

```go
ctx := context.TODO()
id := webpubsub.NewCustomCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue", "customCertificateValue")

read, err := client.CustomCertificatesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebPubSubClient.CustomCertificatesList`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue")

// alternatively `client.CustomCertificatesList(ctx, id)` can be used to do batched pagination
items, err := client.CustomCertificatesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebPubSubClient.CustomDomainsCreateOrUpdate`

```go
ctx := context.TODO()
id := webpubsub.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue", "customDomainValue")

payload := webpubsub.CustomDomain{
	// ...
}


if err := client.CustomDomainsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.CustomDomainsDelete`

```go
ctx := context.TODO()
id := webpubsub.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue", "customDomainValue")

if err := client.CustomDomainsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.CustomDomainsGet`

```go
ctx := context.TODO()
id := webpubsub.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue", "customDomainValue")

read, err := client.CustomDomainsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebPubSubClient.CustomDomainsList`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue")

// alternatively `client.CustomDomainsList(ctx, id)` can be used to do batched pagination
items, err := client.CustomDomainsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebPubSubClient.Delete`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.Get`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebPubSubClient.HubsCreateOrUpdate`

```go
ctx := context.TODO()
id := webpubsub.NewHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue", "hubValue")

payload := webpubsub.WebPubSubHub{
	// ...
}


if err := client.HubsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.HubsDelete`

```go
ctx := context.TODO()
id := webpubsub.NewHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue", "hubValue")

if err := client.HubsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.HubsGet`

```go
ctx := context.TODO()
id := webpubsub.NewHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue", "hubValue")

read, err := client.HubsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebPubSubClient.HubsList`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue")

// alternatively `client.HubsList(ctx, id)` can be used to do batched pagination
items, err := client.HubsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebPubSubClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := webpubsub.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebPubSubClient.ListBySubscription`

```go
ctx := context.TODO()
id := webpubsub.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebPubSubClient.ListKeys`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue")

read, err := client.ListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebPubSubClient.ListSkus`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue")

read, err := client.ListSkus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebPubSubClient.PrivateEndpointConnectionsDelete`

```go
ctx := context.TODO()
id := webpubsub.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue", "privateEndpointConnectionValue")

if err := client.PrivateEndpointConnectionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.PrivateEndpointConnectionsGet`

```go
ctx := context.TODO()
id := webpubsub.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue", "privateEndpointConnectionValue")

read, err := client.PrivateEndpointConnectionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebPubSubClient.PrivateEndpointConnectionsList`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue")

// alternatively `client.PrivateEndpointConnectionsList(ctx, id)` can be used to do batched pagination
items, err := client.PrivateEndpointConnectionsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebPubSubClient.PrivateEndpointConnectionsUpdate`

```go
ctx := context.TODO()
id := webpubsub.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue", "privateEndpointConnectionValue")

payload := webpubsub.PrivateEndpointConnection{
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


### Example Usage: `WebPubSubClient.PrivateLinkResourcesList`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue")

// alternatively `client.PrivateLinkResourcesList(ctx, id)` can be used to do batched pagination
items, err := client.PrivateLinkResourcesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebPubSubClient.RegenerateKey`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue")

payload := webpubsub.RegenerateKeyParameters{
	// ...
}


if err := client.RegenerateKeyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.Restart`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue")

if err := client.RestartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.SharedPrivateLinkResourcesCreateOrUpdate`

```go
ctx := context.TODO()
id := webpubsub.NewSharedPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue", "sharedPrivateLinkResourceValue")

payload := webpubsub.SharedPrivateLinkResource{
	// ...
}


if err := client.SharedPrivateLinkResourcesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.SharedPrivateLinkResourcesDelete`

```go
ctx := context.TODO()
id := webpubsub.NewSharedPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue", "sharedPrivateLinkResourceValue")

if err := client.SharedPrivateLinkResourcesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.SharedPrivateLinkResourcesGet`

```go
ctx := context.TODO()
id := webpubsub.NewSharedPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue", "sharedPrivateLinkResourceValue")

read, err := client.SharedPrivateLinkResourcesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebPubSubClient.SharedPrivateLinkResourcesList`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue")

// alternatively `client.SharedPrivateLinkResourcesList(ctx, id)` can be used to do batched pagination
items, err := client.SharedPrivateLinkResourcesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebPubSubClient.Update`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubValue")

payload := webpubsub.WebPubSubResource{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.UsagesList`

```go
ctx := context.TODO()
id := webpubsub.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

// alternatively `client.UsagesList(ctx, id)` can be used to do batched pagination
items, err := client.UsagesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
