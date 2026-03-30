
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/globalrulestackresources` Documentation

The `globalrulestackresources` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2025-10-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/globalrulestackresources"
```


### Client Initialization

```go
client := globalrulestackresources.NewGlobalRulestackResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GlobalRulestackResourcesClient.GlobalRulestackCreateOrUpdate`

```go
ctx := context.TODO()
id := globalrulestackresources.NewGlobalRulestackID("globalRulestackName")

payload := globalrulestackresources.GlobalRulestackResource{
	// ...
}


if err := client.GlobalRulestackCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GlobalRulestackResourcesClient.GlobalRulestackDelete`

```go
ctx := context.TODO()
id := globalrulestackresources.NewGlobalRulestackID("globalRulestackName")

if err := client.GlobalRulestackDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GlobalRulestackResourcesClient.GlobalRulestackGet`

```go
ctx := context.TODO()
id := globalrulestackresources.NewGlobalRulestackID("globalRulestackName")

read, err := client.GlobalRulestackGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GlobalRulestackResourcesClient.GlobalRulestackList`

```go
ctx := context.TODO()


// alternatively `client.GlobalRulestackList(ctx)` can be used to do batched pagination
items, err := client.GlobalRulestackListComplete(ctx)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalRulestackResourcesClient.GlobalRulestackUpdate`

```go
ctx := context.TODO()
id := globalrulestackresources.NewGlobalRulestackID("globalRulestackName")

payload := globalrulestackresources.GlobalRulestackResourceUpdate{
	// ...
}


read, err := client.GlobalRulestackUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GlobalRulestackResourcesClient.GlobalRulestackcommit`

```go
ctx := context.TODO()
id := globalrulestackresources.NewGlobalRulestackID("globalRulestackName")

if err := client.GlobalRulestackcommitThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GlobalRulestackResourcesClient.GlobalRulestackgetChangeLog`

```go
ctx := context.TODO()
id := globalrulestackresources.NewGlobalRulestackID("globalRulestackName")

read, err := client.GlobalRulestackgetChangeLog(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GlobalRulestackResourcesClient.GlobalRulestacklistAdvancedSecurityObjects`

```go
ctx := context.TODO()
id := globalrulestackresources.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.GlobalRulestacklistAdvancedSecurityObjects(ctx, id, globalrulestackresources.DefaultGlobalRulestacklistAdvancedSecurityObjectsOperationOptions())` can be used to do batched pagination
items, err := client.GlobalRulestacklistAdvancedSecurityObjectsComplete(ctx, id, globalrulestackresources.DefaultGlobalRulestacklistAdvancedSecurityObjectsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalRulestackResourcesClient.GlobalRulestacklistAppIds`

```go
ctx := context.TODO()
id := globalrulestackresources.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.GlobalRulestacklistAppIds(ctx, id, globalrulestackresources.DefaultGlobalRulestacklistAppIdsOperationOptions())` can be used to do batched pagination
items, err := client.GlobalRulestacklistAppIdsComplete(ctx, id, globalrulestackresources.DefaultGlobalRulestacklistAppIdsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalRulestackResourcesClient.GlobalRulestacklistCountries`

```go
ctx := context.TODO()
id := globalrulestackresources.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.GlobalRulestacklistCountries(ctx, id, globalrulestackresources.DefaultGlobalRulestacklistCountriesOperationOptions())` can be used to do batched pagination
items, err := client.GlobalRulestacklistCountriesComplete(ctx, id, globalrulestackresources.DefaultGlobalRulestacklistCountriesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalRulestackResourcesClient.GlobalRulestacklistFirewalls`

```go
ctx := context.TODO()
id := globalrulestackresources.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.GlobalRulestacklistFirewalls(ctx, id)` can be used to do batched pagination
items, err := client.GlobalRulestacklistFirewallsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalRulestackResourcesClient.GlobalRulestacklistPredefinedURLCategories`

```go
ctx := context.TODO()
id := globalrulestackresources.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.GlobalRulestacklistPredefinedURLCategories(ctx, id, globalrulestackresources.DefaultGlobalRulestacklistPredefinedURLCategoriesOperationOptions())` can be used to do batched pagination
items, err := client.GlobalRulestacklistPredefinedURLCategoriesComplete(ctx, id, globalrulestackresources.DefaultGlobalRulestacklistPredefinedURLCategoriesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalRulestackResourcesClient.GlobalRulestacklistSecurityServices`

```go
ctx := context.TODO()
id := globalrulestackresources.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.GlobalRulestacklistSecurityServices(ctx, id, globalrulestackresources.DefaultGlobalRulestacklistSecurityServicesOperationOptions())` can be used to do batched pagination
items, err := client.GlobalRulestacklistSecurityServicesComplete(ctx, id, globalrulestackresources.DefaultGlobalRulestacklistSecurityServicesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalRulestackResourcesClient.GlobalRulestackrevert`

```go
ctx := context.TODO()
id := globalrulestackresources.NewGlobalRulestackID("globalRulestackName")

read, err := client.GlobalRulestackrevert(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
