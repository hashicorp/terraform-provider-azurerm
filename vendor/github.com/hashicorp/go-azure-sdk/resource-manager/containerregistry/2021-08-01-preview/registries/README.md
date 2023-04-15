
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/registries` Documentation

The `registries` SDK allows for interaction with the Azure Resource Manager Service `containerregistry` (API Version `2021-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/registries"
```


### Client Initialization

```go
client := registries.NewRegistriesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RegistriesClient.Create`

```go
ctx := context.TODO()
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

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
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RegistriesClient.GenerateCredentials`

```go
ctx := context.TODO()
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

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
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

read, err := client.Get(ctx, id)
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
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

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
id := registries.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

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
id := registries.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

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
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

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
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

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
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

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
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

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
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

payload := registries.RegistryUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
