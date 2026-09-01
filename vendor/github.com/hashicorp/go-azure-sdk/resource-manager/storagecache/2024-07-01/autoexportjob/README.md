
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/autoexportjob` Documentation

The `autoexportjob` SDK allows for interaction with Azure Resource Manager `storagecache` (API Version `2024-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/autoexportjob"
```


### Client Initialization

```go
client := autoexportjob.NewAutoExportJobClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AutoExportJobClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := autoexportjob.NewAutoExportJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName", "autoExportJobName")

payload := autoexportjob.AutoExportJob{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AutoExportJobClient.ListByAmlFilesystem`

```go
ctx := context.TODO()
id := autoexportjob.NewAmlFilesystemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName")

// alternatively `client.ListByAmlFilesystem(ctx, id)` can be used to do batched pagination
items, err := client.ListByAmlFilesystemComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AutoExportJobClient.Update`

```go
ctx := context.TODO()
id := autoexportjob.NewAutoExportJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "amlFilesystemName", "autoExportJobName")

payload := autoexportjob.AutoExportJobUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
