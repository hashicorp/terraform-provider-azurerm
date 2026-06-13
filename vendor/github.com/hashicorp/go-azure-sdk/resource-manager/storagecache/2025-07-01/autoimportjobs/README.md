
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2025-07-01/autoimportjobs` Documentation

The `autoimportjobs` SDK allows for interaction with Azure Resource Manager `storagecache` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2025-07-01/autoimportjobs"
```


### Client Initialization

```go
client := autoimportjobs.NewAutoImportJobsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AutoImportJobsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := autoimportjobs.NewAutoImportJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName", "autoImportJobName")

payload := autoimportjobs.AutoImportJob{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AutoImportJobsClient.Delete`

```go
ctx := context.TODO()
id := autoimportjobs.NewAutoImportJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName", "autoImportJobName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AutoImportJobsClient.Get`

```go
ctx := context.TODO()
id := autoimportjobs.NewAutoImportJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName", "autoImportJobName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AutoImportJobsClient.ListByAmlFilesystem`

```go
ctx := context.TODO()
id := autoimportjobs.NewAmlFilesystemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName")

// alternatively `client.ListByAmlFilesystem(ctx, id)` can be used to do batched pagination
items, err := client.ListByAmlFilesystemComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AutoImportJobsClient.Update`

```go
ctx := context.TODO()
id := autoimportjobs.NewAutoImportJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName", "autoImportJobName")

payload := autoimportjobs.AutoImportJobUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
