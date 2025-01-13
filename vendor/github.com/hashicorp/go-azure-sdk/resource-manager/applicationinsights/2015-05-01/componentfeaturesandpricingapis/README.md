
## `github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentfeaturesandpricingapis` Documentation

The `componentfeaturesandpricingapis` SDK allows for interaction with Azure Resource Manager `applicationinsights` (API Version `2015-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentfeaturesandpricingapis"
```


### Client Initialization

```go
client := componentfeaturesandpricingapis.NewComponentFeaturesAndPricingAPIsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ComponentFeaturesAndPricingAPIsClient.ComponentAvailableFeaturesGet`

```go
ctx := context.TODO()
id := componentfeaturesandpricingapis.NewComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName")

read, err := client.ComponentAvailableFeaturesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentFeaturesAndPricingAPIsClient.ComponentCurrentBillingFeaturesGet`

```go
ctx := context.TODO()
id := componentfeaturesandpricingapis.NewComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName")

read, err := client.ComponentCurrentBillingFeaturesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentFeaturesAndPricingAPIsClient.ComponentCurrentBillingFeaturesUpdate`

```go
ctx := context.TODO()
id := componentfeaturesandpricingapis.NewComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName")

payload := componentfeaturesandpricingapis.ApplicationInsightsComponentBillingFeatures{
	// ...
}


read, err := client.ComponentCurrentBillingFeaturesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentFeaturesAndPricingAPIsClient.ComponentFeatureCapabilitiesGet`

```go
ctx := context.TODO()
id := componentfeaturesandpricingapis.NewComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName")

read, err := client.ComponentFeatureCapabilitiesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentFeaturesAndPricingAPIsClient.ComponentQuotaStatusGet`

```go
ctx := context.TODO()
id := componentfeaturesandpricingapis.NewComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName")

read, err := client.ComponentQuotaStatusGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
