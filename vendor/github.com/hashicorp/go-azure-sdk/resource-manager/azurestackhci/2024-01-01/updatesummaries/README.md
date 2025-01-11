
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/updatesummaries` Documentation

The `updatesummaries` SDK allows for interaction with Azure Resource Manager `azurestackhci` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/updatesummaries"
```


### Client Initialization

```go
client := updatesummaries.NewUpdateSummariesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `UpdateSummariesClient.Delete`

```go
ctx := context.TODO()
id := updatesummaries.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `UpdateSummariesClient.Get`

```go
ctx := context.TODO()
id := updatesummaries.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `UpdateSummariesClient.List`

```go
ctx := context.TODO()
id := updatesummaries.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `UpdateSummariesClient.Put`

```go
ctx := context.TODO()
id := updatesummaries.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := updatesummaries.UpdateSummaries{
	// ...
}


read, err := client.Put(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
