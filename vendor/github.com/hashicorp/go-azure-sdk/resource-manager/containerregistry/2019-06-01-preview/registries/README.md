
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/registries` Documentation

The `registries` SDK allows for interaction with Azure Resource Manager `containerregistry` (API Version `2019-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/registries"
```


### Client Initialization

```go
client := registries.NewRegistriesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RegistriesClient.GetBuildSourceUploadURL`

```go
ctx := context.TODO()
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

read, err := client.GetBuildSourceUploadURL(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegistriesClient.ScheduleRun`

```go
ctx := context.TODO()
id := registries.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

payload := registries.RunRequest{
	// ...
}


if err := client.ScheduleRunThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
