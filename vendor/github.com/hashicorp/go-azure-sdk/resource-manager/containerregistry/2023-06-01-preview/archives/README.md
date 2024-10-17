
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/archives` Documentation

The `archives` SDK allows for interaction with Azure Resource Manager `containerregistry` (API Version `2023-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/archives"
```


### Client Initialization

```go
client := archives.NewArchivesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ArchivesClient.Create`

```go
ctx := context.TODO()
id := archives.NewArchiveID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "packageName", "archiveName")

payload := archives.Archive{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ArchivesClient.Delete`

```go
ctx := context.TODO()
id := archives.NewArchiveID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "packageName", "archiveName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ArchivesClient.Get`

```go
ctx := context.TODO()
id := archives.NewArchiveID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "packageName", "archiveName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ArchivesClient.List`

```go
ctx := context.TODO()
id := archives.NewPackageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "packageName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ArchivesClient.Update`

```go
ctx := context.TODO()
id := archives.NewArchiveID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "packageName", "archiveName")

payload := archives.ArchiveUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
