
## `github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-07-01/encodings` Documentation

The `encodings` SDK allows for interaction with the Azure Resource Manager Service `media` (API Version `2022-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-07-01/encodings"
```


### Client Initialization

```go
client := encodings.NewEncodingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EncodingsClient.JobsCancelJob`

```go
ctx := context.TODO()
id := encodings.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "transformValue", "jobValue")

read, err := client.JobsCancelJob(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EncodingsClient.JobsCreate`

```go
ctx := context.TODO()
id := encodings.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "transformValue", "jobValue")

payload := encodings.Job{
	// ...
}


read, err := client.JobsCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EncodingsClient.JobsDelete`

```go
ctx := context.TODO()
id := encodings.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "transformValue", "jobValue")

read, err := client.JobsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EncodingsClient.JobsGet`

```go
ctx := context.TODO()
id := encodings.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "transformValue", "jobValue")

read, err := client.JobsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EncodingsClient.JobsList`

```go
ctx := context.TODO()
id := encodings.NewTransformID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "transformValue")

// alternatively `client.JobsList(ctx, id, encodings.DefaultJobsListOperationOptions())` can be used to do batched pagination
items, err := client.JobsListComplete(ctx, id, encodings.DefaultJobsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EncodingsClient.JobsUpdate`

```go
ctx := context.TODO()
id := encodings.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "transformValue", "jobValue")

payload := encodings.Job{
	// ...
}


read, err := client.JobsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EncodingsClient.TransformsCreateOrUpdate`

```go
ctx := context.TODO()
id := encodings.NewTransformID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "transformValue")

payload := encodings.Transform{
	// ...
}


read, err := client.TransformsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EncodingsClient.TransformsDelete`

```go
ctx := context.TODO()
id := encodings.NewTransformID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "transformValue")

read, err := client.TransformsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EncodingsClient.TransformsGet`

```go
ctx := context.TODO()
id := encodings.NewTransformID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "transformValue")

read, err := client.TransformsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EncodingsClient.TransformsList`

```go
ctx := context.TODO()
id := encodings.NewMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

// alternatively `client.TransformsList(ctx, id, encodings.DefaultTransformsListOperationOptions())` can be used to do batched pagination
items, err := client.TransformsListComplete(ctx, id, encodings.DefaultTransformsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EncodingsClient.TransformsUpdate`

```go
ctx := context.TODO()
id := encodings.NewTransformID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "transformValue")

payload := encodings.Transform{
	// ...
}


read, err := client.TransformsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
