
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/importjobs` Documentation

The `importjobs` SDK allows for interaction with Azure Resource Manager `storagecache` (API Version `2024-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/importjobs"
```


### Client Initialization

```go
client := importjobs.NewImportJobsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ImportJobsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := importjobs.NewImportJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName", "importJobName")

payload := importjobs.ImportJob{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ImportJobsClient.Delete`

```go
ctx := context.TODO()
id := importjobs.NewImportJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName", "importJobName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ImportJobsClient.Get`

```go
ctx := context.TODO()
id := importjobs.NewImportJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName", "importJobName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ImportJobsClient.ListByAmlFilesystem`

```go
ctx := context.TODO()
id := importjobs.NewAmlFilesystemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName")

// alternatively `client.ListByAmlFilesystem(ctx, id)` can be used to do batched pagination
items, err := client.ListByAmlFilesystemComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ImportJobsClient.Update`

```go
ctx := context.TODO()
id := importjobs.NewImportJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName", "importJobName")

payload := importjobs.ImportJobUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
