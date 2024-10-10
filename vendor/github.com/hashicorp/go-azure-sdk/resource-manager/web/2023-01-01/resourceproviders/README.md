
## `github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders` Documentation

The `resourceproviders` SDK allows for interaction with Azure Resource Manager `web` (API Version `2023-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
```


### Client Initialization

```go
client := resourceproviders.NewResourceProvidersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ResourceProvidersClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := resourceproviders.ResourceNameAvailabilityRequest{
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


### Example Usage: `ResourceProvidersClient.GetPublishingUser`

```go
ctx := context.TODO()


read, err := client.GetPublishingUser(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceProvidersClient.GetSourceControl`

```go
ctx := context.TODO()
id := resourceproviders.NewSourceControlID("sourceControlName")

read, err := client.GetSourceControl(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceProvidersClient.GetSubscriptionDeploymentLocations`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.GetSubscriptionDeploymentLocations(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceProvidersClient.GetUsagesInLocationlist`

```go
ctx := context.TODO()
id := resourceproviders.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.GetUsagesInLocationlist(ctx, id)` can be used to do batched pagination
items, err := client.GetUsagesInLocationlistComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceProvidersClient.ListAseRegions`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAseRegions(ctx, id)` can be used to do batched pagination
items, err := client.ListAseRegionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceProvidersClient.ListBillingMeters`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBillingMeters(ctx, id, resourceproviders.DefaultListBillingMetersOperationOptions())` can be used to do batched pagination
items, err := client.ListBillingMetersComplete(ctx, id, resourceproviders.DefaultListBillingMetersOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceProvidersClient.ListCustomHostNameSites`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListCustomHostNameSites(ctx, id, resourceproviders.DefaultListCustomHostNameSitesOperationOptions())` can be used to do batched pagination
items, err := client.ListCustomHostNameSitesComplete(ctx, id, resourceproviders.DefaultListCustomHostNameSitesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceProvidersClient.ListGeoRegions`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListGeoRegions(ctx, id, resourceproviders.DefaultListGeoRegionsOperationOptions())` can be used to do batched pagination
items, err := client.ListGeoRegionsComplete(ctx, id, resourceproviders.DefaultListGeoRegionsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceProvidersClient.ListPremierAddOnOffers`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListPremierAddOnOffers(ctx, id)` can be used to do batched pagination
items, err := client.ListPremierAddOnOffersComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceProvidersClient.ListSiteIdentifiersAssignedToHostName`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := resourceproviders.NameIdentifier{
	// ...
}


// alternatively `client.ListSiteIdentifiersAssignedToHostName(ctx, id, payload)` can be used to do batched pagination
items, err := client.ListSiteIdentifiersAssignedToHostNameComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceProvidersClient.ListSkus`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListSkus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceProvidersClient.ListSourceControls`

```go
ctx := context.TODO()


// alternatively `client.ListSourceControls(ctx)` can be used to do batched pagination
items, err := client.ListSourceControlsComplete(ctx)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceProvidersClient.Move`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := resourceproviders.CsmMoveResourceEnvelope{
	// ...
}


read, err := client.Move(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceProvidersClient.UpdatePublishingUser`

```go
ctx := context.TODO()

payload := resourceproviders.User{
	// ...
}


read, err := client.UpdatePublishingUser(ctx, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceProvidersClient.UpdateSourceControl`

```go
ctx := context.TODO()
id := resourceproviders.NewSourceControlID("sourceControlName")

payload := resourceproviders.SourceControl{
	// ...
}


read, err := client.UpdateSourceControl(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceProvidersClient.Validate`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := resourceproviders.ValidateRequest{
	// ...
}


read, err := client.Validate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceProvidersClient.ValidateMove`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := resourceproviders.CsmMoveResourceEnvelope{
	// ...
}


read, err := client.ValidateMove(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceProvidersClient.VerifyHostingEnvironmentVnet`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := resourceproviders.VnetParameters{
	// ...
}


read, err := client.VerifyHostingEnvironmentVnet(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
