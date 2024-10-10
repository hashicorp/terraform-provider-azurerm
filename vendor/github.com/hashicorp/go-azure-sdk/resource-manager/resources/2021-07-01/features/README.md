
## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features` Documentation

The `features` SDK allows for interaction with Azure Resource Manager `resources` (API Version `2021-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features"
```


### Client Initialization

```go
client := features.NewFeaturesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FeaturesClient.Get`

```go
ctx := context.TODO()
id := features.NewFeatureID("12345678-1234-9876-4563-123456789012", "providerName", "featureName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FeaturesClient.List`

```go
ctx := context.TODO()
id := features.NewProviders2ID("12345678-1234-9876-4563-123456789012", "providerName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FeaturesClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FeaturesClient.Register`

```go
ctx := context.TODO()
id := features.NewFeatureID("12345678-1234-9876-4563-123456789012", "providerName", "featureName")

read, err := client.Register(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FeaturesClient.Unregister`

```go
ctx := context.TODO()
id := features.NewFeatureID("12345678-1234-9876-4563-123456789012", "providerName", "featureName")

read, err := client.Unregister(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
