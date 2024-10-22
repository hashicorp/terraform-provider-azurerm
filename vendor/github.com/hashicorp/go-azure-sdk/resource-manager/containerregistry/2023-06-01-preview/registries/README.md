
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/registries` Documentation

The `registries` SDK allows for interaction with Azure Resource Manager `containerregistry` (API Version `2023-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/registries"
```


### Client Initialization

```go
client := registries.NewRegistriesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RegistriesClient.Create`

```go
ctx := context.TODO()
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

payload := registries.Registry{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RegistriesClient.Delete`

```go
ctx := context.TODO()
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RegistriesClient.GenerateCredentials`

```go
ctx := context.TODO()
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

payload := registries.GenerateCredentialsParameters{
	// ...
}


if err := client.GenerateCredentialsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RegistriesClient.Get`

```go
ctx := context.TODO()
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegistriesClient.GetPrivateLinkResource`

```go
ctx := context.TODO()
id := registries.NewPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "privateLinkResourceName")

read, err := client.GetPrivateLinkResource(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegistriesClient.ImportImage`

```go
ctx := context.TODO()
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

payload := registries.ImportImageParameters{
	// ...
}


if err := client.ImportImageThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RegistriesClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RegistriesClient.ListByResourceGroup`

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


### Example Usage: `RegistriesClient.ListCredentials`

```go
ctx := context.TODO()
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

read, err := client.ListCredentials(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegistriesClient.ListPrivateLinkResources`

```go
ctx := context.TODO()
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

// alternatively `client.ListPrivateLinkResources(ctx, id)` can be used to do batched pagination
items, err := client.ListPrivateLinkResourcesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RegistriesClient.ListUsages`

```go
ctx := context.TODO()
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

read, err := client.ListUsages(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegistriesClient.RegenerateCredential`

```go
ctx := context.TODO()
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

payload := registries.RegenerateCredentialParameters{
	// ...
}


read, err := client.RegenerateCredential(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegistriesClient.Update`

```go
ctx := context.TODO()
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

payload := registries.RegistryUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
