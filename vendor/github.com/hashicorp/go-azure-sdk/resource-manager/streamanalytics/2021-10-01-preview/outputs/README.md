
## `github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs` Documentation

The `outputs` SDK allows for interaction with the Azure Resource Manager Service `streamanalytics` (API Version `2021-10-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
```


### Client Initialization

```go
client := outputs.NewOutputsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OutputsClient.CreateOrReplace`

```go
ctx := context.TODO()
id := outputs.NewOutputID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue", "outputValue")

payload := outputs.Output{
	// ...
}


read, err := client.CreateOrReplace(ctx, id, payload, outputs.DefaultCreateOrReplaceOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OutputsClient.Delete`

```go
ctx := context.TODO()
id := outputs.NewOutputID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue", "outputValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OutputsClient.Get`

```go
ctx := context.TODO()
id := outputs.NewOutputID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue", "outputValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OutputsClient.ListByStreamingJob`

```go
ctx := context.TODO()
id := outputs.NewStreamingJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue")

// alternatively `client.ListByStreamingJob(ctx, id, outputs.DefaultListByStreamingJobOperationOptions())` can be used to do batched pagination
items, err := client.ListByStreamingJobComplete(ctx, id, outputs.DefaultListByStreamingJobOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OutputsClient.Test`

```go
ctx := context.TODO()
id := outputs.NewOutputID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue", "outputValue")

payload := outputs.Output{
	// ...
}


if err := client.TestThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OutputsClient.Update`

```go
ctx := context.TODO()
id := outputs.NewOutputID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue", "outputValue")

payload := outputs.Output{
	// ...
}


read, err := client.Update(ctx, id, payload, outputs.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
