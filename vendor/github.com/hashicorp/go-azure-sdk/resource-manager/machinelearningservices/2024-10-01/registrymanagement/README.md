
## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-10-01/registrymanagement` Documentation

The `registrymanagement` SDK allows for interaction with Azure Resource Manager `machinelearningservices` (API Version `2024-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-10-01/registrymanagement"
```


### Client Initialization

```go
client := registrymanagement.NewRegistryManagementClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RegistryManagementClient.RegistriesCreateOrUpdate`

```go
ctx := context.TODO()
id := registrymanagement.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

payload := registrymanagement.RegistryTrackedResource{
	// ...
}


if err := client.RegistriesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RegistryManagementClient.RegistriesDelete`

```go
ctx := context.TODO()
id := registrymanagement.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

if err := client.RegistriesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RegistryManagementClient.RegistriesGet`

```go
ctx := context.TODO()
id := registrymanagement.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

read, err := client.RegistriesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegistryManagementClient.RegistriesList`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.RegistriesList(ctx, id)` can be used to do batched pagination
items, err := client.RegistriesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RegistryManagementClient.RegistriesListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.RegistriesListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.RegistriesListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RegistryManagementClient.RegistriesRemoveRegions`

```go
ctx := context.TODO()
id := registrymanagement.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

payload := registrymanagement.RegistryTrackedResource{
	// ...
}


if err := client.RegistriesRemoveRegionsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RegistryManagementClient.RegistriesUpdate`

```go
ctx := context.TODO()
id := registrymanagement.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

payload := registrymanagement.PartialRegistryPartialTrackedResource{
	// ...
}


read, err := client.RegistriesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
