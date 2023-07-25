
## `github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/transformations` Documentation

The `transformations` SDK allows for interaction with the Azure Resource Manager Service `streamanalytics` (API Version `2020-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/transformations"
```


### Client Initialization

```go
client := transformations.NewTransformationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TransformationsClient.CreateOrReplace`

```go
ctx := context.TODO()
id := transformations.NewTransformationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue", "transformationValue")

payload := transformations.Transformation{
	// ...
}


read, err := client.CreateOrReplace(ctx, id, payload, transformations.DefaultCreateOrReplaceOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TransformationsClient.Get`

```go
ctx := context.TODO()
id := transformations.NewTransformationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue", "transformationValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TransformationsClient.Update`

```go
ctx := context.TODO()
id := transformations.NewTransformationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "streamingJobValue", "transformationValue")

payload := transformations.Transformation{
	// ...
}


read, err := client.Update(ctx, id, payload, transformations.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
