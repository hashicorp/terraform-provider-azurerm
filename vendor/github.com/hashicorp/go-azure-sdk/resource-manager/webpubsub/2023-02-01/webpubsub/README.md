
## `github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2023-02-01/webpubsub` Documentation

The `webpubsub` SDK allows for interaction with Azure Resource Manager `webpubsub` (API Version `2023-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
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
id := webpubsub.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

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
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName")

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
id := webpubsub.NewCustomCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName", "customCertificateName")

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
id := webpubsub.NewCustomCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName", "customCertificateName")

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
id := webpubsub.NewCustomCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName", "customCertificateName")

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
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName")

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
id := webpubsub.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName", "customDomainName")

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
id := webpubsub.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName", "customDomainName")

if err := client.CustomDomainsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.CustomDomainsGet`

```go
ctx := context.TODO()
id := webpubsub.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName", "customDomainName")

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
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName")

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
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.Get`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName")

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
id := webpubsub.NewHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName", "hubName")

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
id := webpubsub.NewHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName", "hubName")

if err := client.HubsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.HubsGet`

```go
ctx := context.TODO()
id := webpubsub.NewHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName", "hubName")

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
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName")

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


### Example Usage: `WebPubSubClient.ListBySubscription`

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


### Example Usage: `WebPubSubClient.ListKeys`

```go
ctx := context.TODO()
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName")

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
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName")

// alternatively `client.ListSkus(ctx, id)` can be used to do batched pagination
items, err := client.ListSkusComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebPubSubClient.PrivateEndpointConnectionsDelete`

```go
ctx := context.TODO()
id := webpubsub.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName", "privateEndpointConnectionName")

if err := client.PrivateEndpointConnectionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.PrivateEndpointConnectionsGet`

```go
ctx := context.TODO()
id := webpubsub.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName", "privateEndpointConnectionName")

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
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName")

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
id := webpubsub.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName", "privateEndpointConnectionName")

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
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName")

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
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName")

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
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName")

if err := client.RestartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.SharedPrivateLinkResourcesCreateOrUpdate`

```go
ctx := context.TODO()
id := webpubsub.NewSharedPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName", "sharedPrivateLinkResourceName")

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
id := webpubsub.NewSharedPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName", "sharedPrivateLinkResourceName")

if err := client.SharedPrivateLinkResourcesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebPubSubClient.SharedPrivateLinkResourcesGet`

```go
ctx := context.TODO()
id := webpubsub.NewSharedPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName", "sharedPrivateLinkResourceName")

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
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName")

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
id := webpubsub.NewWebPubSubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webPubSubName")

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
id := webpubsub.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.UsagesList(ctx, id)` can be used to do batched pagination
items, err := client.UsagesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
