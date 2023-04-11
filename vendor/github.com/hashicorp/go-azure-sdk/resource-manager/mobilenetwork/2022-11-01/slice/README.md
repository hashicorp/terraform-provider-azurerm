
## `github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/slice` Documentation

The `slice` SDK allows for interaction with the Azure Resource Manager Service `mobilenetwork` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/slice"
```


### Client Initialization

```go
client := slice.NewSliceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SliceClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := slice.NewSliceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkValue", "sliceValue")

payload := slice.Slice{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SliceClient.Delete`

```go
ctx := context.TODO()
id := slice.NewSliceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkValue", "sliceValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SliceClient.Get`

```go
ctx := context.TODO()
id := slice.NewSliceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkValue", "sliceValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SliceClient.UpdateTags`

```go
ctx := context.TODO()
id := slice.NewSliceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkValue", "sliceValue")

payload := slice.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
