
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/archiveversions` Documentation

The `archiveversions` SDK allows for interaction with Azure Resource Manager `containerregistry` (API Version `2023-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/archiveversions"
```


### Client Initialization

```go
client := archiveversions.NewArchiveVersionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ArchiveVersionsClient.Create`

```go
ctx := context.TODO()
id := archiveversions.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "packageName", "archiveName", "versionName")

if err := client.CreateThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ArchiveVersionsClient.Delete`

```go
ctx := context.TODO()
id := archiveversions.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "packageName", "archiveName", "versionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ArchiveVersionsClient.Get`

```go
ctx := context.TODO()
id := archiveversions.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "packageName", "archiveName", "versionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ArchiveVersionsClient.List`

```go
ctx := context.TODO()
id := archiveversions.NewArchiveID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "packageName", "archiveName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
