
## `github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/functions` Documentation

The `functions` SDK allows for interaction with the Azure Resource Manager Service `streamanalytics` (API Version `2020-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/functions"
```


### Client Initialization

```go
client := functions.NewFunctionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FunctionsClient.CreateOrReplace`

```go
ctx := context.TODO()
id := functions.NewFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue", "functionValue")

payload := functions.Function{
	// ...
}


read, err := client.CreateOrReplace(ctx, id, payload, functions.DefaultCreateOrReplaceOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FunctionsClient.Delete`

```go
ctx := context.TODO()
id := functions.NewFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue", "functionValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FunctionsClient.Get`

```go
ctx := context.TODO()
id := functions.NewFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue", "functionValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FunctionsClient.ListByStreamingJob`

```go
ctx := context.TODO()
id := functions.NewStreamingJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue")

// alternatively `client.ListByStreamingJob(ctx, id, functions.DefaultListByStreamingJobOperationOptions())` can be used to do batched pagination
items, err := client.ListByStreamingJobComplete(ctx, id, functions.DefaultListByStreamingJobOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FunctionsClient.RetrieveDefaultDefinition`

```go
ctx := context.TODO()
id := functions.NewFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue", "functionValue")

payload := functions.FunctionRetrieveDefaultDefinitionParameters{
	// ...
}


read, err := client.RetrieveDefaultDefinition(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FunctionsClient.Test`

```go
ctx := context.TODO()
id := functions.NewFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue", "functionValue")

payload := functions.Function{
	// ...
}


if err := client.TestThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FunctionsClient.Update`

```go
ctx := context.TODO()
id := functions.NewFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue", "functionValue")

payload := functions.Function{
	// ...
}


read, err := client.Update(ctx, id, payload, functions.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
