
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources` Documentation

The `localrulestackresources` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2025-10-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
```


### Client Initialization

```go
client := localrulestackresources.NewLocalRulestackResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestacksCreateOrUpdate`

```go
ctx := context.TODO()
id := localrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

payload := localrulestackresources.LocalRulestackResource{
	// ...
}


if err := client.LocalRulestacksCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestacksDelete`

```go
ctx := context.TODO()
id := localrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

if err := client.LocalRulestacksDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestacksGet`

```go
ctx := context.TODO()
id := localrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

read, err := client.LocalRulestacksGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestacksListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.LocalRulestacksListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.LocalRulestacksListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestacksListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.LocalRulestacksListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.LocalRulestacksListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestacksUpdate`

```go
ctx := context.TODO()
id := localrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

payload := localrulestackresources.LocalRulestackResourceUpdate{
	// ...
}


read, err := client.LocalRulestacksUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestackscommit`

```go
ctx := context.TODO()
id := localrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

if err := client.LocalRulestackscommitThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestacksgetChangeLog`

```go
ctx := context.TODO()
id := localrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

read, err := client.LocalRulestacksgetChangeLog(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestacksgetSupportInfo`

```go
ctx := context.TODO()
id := localrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

read, err := client.LocalRulestacksgetSupportInfo(ctx, id, localrulestackresources.DefaultLocalRulestacksgetSupportInfoOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestackslistAdvancedSecurityObjects`

```go
ctx := context.TODO()
id := localrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.LocalRulestackslistAdvancedSecurityObjects(ctx, id, localrulestackresources.DefaultLocalRulestackslistAdvancedSecurityObjectsOperationOptions())` can be used to do batched pagination
items, err := client.LocalRulestackslistAdvancedSecurityObjectsComplete(ctx, id, localrulestackresources.DefaultLocalRulestackslistAdvancedSecurityObjectsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestackslistAppIds`

```go
ctx := context.TODO()
id := localrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.LocalRulestackslistAppIds(ctx, id, localrulestackresources.DefaultLocalRulestackslistAppIdsOperationOptions())` can be used to do batched pagination
items, err := client.LocalRulestackslistAppIdsComplete(ctx, id, localrulestackresources.DefaultLocalRulestackslistAppIdsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestackslistCountries`

```go
ctx := context.TODO()
id := localrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.LocalRulestackslistCountries(ctx, id, localrulestackresources.DefaultLocalRulestackslistCountriesOperationOptions())` can be used to do batched pagination
items, err := client.LocalRulestackslistCountriesComplete(ctx, id, localrulestackresources.DefaultLocalRulestackslistCountriesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestackslistFirewalls`

```go
ctx := context.TODO()
id := localrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.LocalRulestackslistFirewalls(ctx, id)` can be used to do batched pagination
items, err := client.LocalRulestackslistFirewallsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestackslistPredefinedURLCategories`

```go
ctx := context.TODO()
id := localrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.LocalRulestackslistPredefinedURLCategories(ctx, id, localrulestackresources.DefaultLocalRulestackslistPredefinedURLCategoriesOperationOptions())` can be used to do batched pagination
items, err := client.LocalRulestackslistPredefinedURLCategoriesComplete(ctx, id, localrulestackresources.DefaultLocalRulestackslistPredefinedURLCategoriesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestackslistSecurityServices`

```go
ctx := context.TODO()
id := localrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.LocalRulestackslistSecurityServices(ctx, id, localrulestackresources.DefaultLocalRulestackslistSecurityServicesOperationOptions())` can be used to do batched pagination
items, err := client.LocalRulestackslistSecurityServicesComplete(ctx, id, localrulestackresources.DefaultLocalRulestackslistSecurityServicesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulestackResourcesClient.LocalRulestacksrevert`

```go
ctx := context.TODO()
id := localrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

read, err := client.LocalRulestacksrevert(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
