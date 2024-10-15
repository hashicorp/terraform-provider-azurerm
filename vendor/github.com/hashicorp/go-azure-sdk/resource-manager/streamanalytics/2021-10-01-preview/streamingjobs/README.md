
## `github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs` Documentation

The `streamingjobs` SDK allows for interaction with Azure Resource Manager `streamanalytics` (API Version `2021-10-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
```


### Client Initialization

```go
client := streamingjobs.NewStreamingJobsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StreamingJobsClient.CreateOrReplace`

```go
ctx := context.TODO()
id := streamingjobs.NewStreamingJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobName")

payload := streamingjobs.StreamingJob{
	// ...
}


if err := client.CreateOrReplaceThenPoll(ctx, id, payload, streamingjobs.DefaultCreateOrReplaceOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `StreamingJobsClient.Delete`

```go
ctx := context.TODO()
id := streamingjobs.NewStreamingJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StreamingJobsClient.Get`

```go
ctx := context.TODO()
id := streamingjobs.NewStreamingJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobName")

read, err := client.Get(ctx, id, streamingjobs.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StreamingJobsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id, streamingjobs.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, streamingjobs.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StreamingJobsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, streamingjobs.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, streamingjobs.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StreamingJobsClient.Scale`

```go
ctx := context.TODO()
id := streamingjobs.NewStreamingJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobName")

payload := streamingjobs.ScaleStreamingJobParameters{
	// ...
}


if err := client.ScaleThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StreamingJobsClient.SkuList`

```go
ctx := context.TODO()
id := streamingjobs.NewStreamingJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobName")

// alternatively `client.SkuList(ctx, id)` can be used to do batched pagination
items, err := client.SkuListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StreamingJobsClient.Start`

```go
ctx := context.TODO()
id := streamingjobs.NewStreamingJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobName")

payload := streamingjobs.StartStreamingJobParameters{
	// ...
}


if err := client.StartThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StreamingJobsClient.Stop`

```go
ctx := context.TODO()
id := streamingjobs.NewStreamingJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobName")

if err := client.StopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StreamingJobsClient.Update`

```go
ctx := context.TODO()
id := streamingjobs.NewStreamingJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobName")

payload := streamingjobs.StreamingJob{
	// ...
}


read, err := client.Update(ctx, id, payload, streamingjobs.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
