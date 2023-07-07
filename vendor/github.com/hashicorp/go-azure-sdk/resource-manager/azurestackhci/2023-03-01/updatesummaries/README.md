
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/updatesummaries` Documentation

The `updatesummaries` SDK allows for interaction with the Azure Resource Manager Service `azurestackhci` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/updatesummaries"
```


### Client Initialization

```go
client := updatesummaries.NewUpdateSummariesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `UpdateSummariesClient.UpdateSummariesDelete`

```go
ctx := context.TODO()
id := updatesummaries.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

if err := client.UpdateSummariesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `UpdateSummariesClient.UpdateSummariesGet`

```go
ctx := context.TODO()
id := updatesummaries.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

read, err := client.UpdateSummariesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `UpdateSummariesClient.UpdateSummariesList`

```go
ctx := context.TODO()
id := updatesummaries.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

// alternatively `client.UpdateSummariesList(ctx, id)` can be used to do batched pagination
items, err := client.UpdateSummariesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `UpdateSummariesClient.UpdateSummariesPut`

```go
ctx := context.TODO()
id := updatesummaries.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := updatesummaries.UpdateSummaries{
	// ...
}


read, err := client.UpdateSummariesPut(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
