
## `github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs` Documentation

The `inputs` SDK allows for interaction with Azure Resource Manager `streamanalytics` (API Version `2020-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
```


### Client Initialization

```go
client := inputs.NewInputsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `InputsClient.CreateOrReplace`

```go
ctx := context.TODO()
id := inputs.NewInputID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobName", "inputName")

payload := inputs.Input{
	// ...
}


read, err := client.CreateOrReplace(ctx, id, payload, inputs.DefaultCreateOrReplaceOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `InputsClient.Delete`

```go
ctx := context.TODO()
id := inputs.NewInputID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobName", "inputName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `InputsClient.Get`

```go
ctx := context.TODO()
id := inputs.NewInputID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobName", "inputName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `InputsClient.ListByStreamingJob`

```go
ctx := context.TODO()
id := inputs.NewStreamingJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobName")

// alternatively `client.ListByStreamingJob(ctx, id, inputs.DefaultListByStreamingJobOperationOptions())` can be used to do batched pagination
items, err := client.ListByStreamingJobComplete(ctx, id, inputs.DefaultListByStreamingJobOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `InputsClient.Test`

```go
ctx := context.TODO()
id := inputs.NewInputID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobName", "inputName")

payload := inputs.Input{
	// ...
}


if err := client.TestThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `InputsClient.Update`

```go
ctx := context.TODO()
id := inputs.NewInputID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobName", "inputName")

payload := inputs.Input{
	// ...
}


read, err := client.Update(ctx, id, payload, inputs.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
