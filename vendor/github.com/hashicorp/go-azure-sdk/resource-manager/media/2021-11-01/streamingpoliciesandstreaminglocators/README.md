
## `github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/streamingpoliciesandstreaminglocators` Documentation

The `streamingpoliciesandstreaminglocators` SDK allows for interaction with the Azure Resource Manager Service `media` (API Version `2021-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/streamingpoliciesandstreaminglocators"
```


### Client Initialization

```go
client := streamingpoliciesandstreaminglocators.NewStreamingPoliciesAndStreamingLocatorsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StreamingPoliciesAndStreamingLocatorsClient.StreamingLocatorsCreate`

```go
ctx := context.TODO()
id := streamingpoliciesandstreaminglocators.NewStreamingLocatorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingLocatorValue")

payload := streamingpoliciesandstreaminglocators.StreamingLocator{
	// ...
}


read, err := client.StreamingLocatorsCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StreamingPoliciesAndStreamingLocatorsClient.StreamingLocatorsDelete`

```go
ctx := context.TODO()
id := streamingpoliciesandstreaminglocators.NewStreamingLocatorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingLocatorValue")

read, err := client.StreamingLocatorsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StreamingPoliciesAndStreamingLocatorsClient.StreamingLocatorsGet`

```go
ctx := context.TODO()
id := streamingpoliciesandstreaminglocators.NewStreamingLocatorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingLocatorValue")

read, err := client.StreamingLocatorsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StreamingPoliciesAndStreamingLocatorsClient.StreamingLocatorsList`

```go
ctx := context.TODO()
id := streamingpoliciesandstreaminglocators.NewMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

// alternatively `client.StreamingLocatorsList(ctx, id, streamingpoliciesandstreaminglocators.DefaultStreamingLocatorsListOperationOptions())` can be used to do batched pagination
items, err := client.StreamingLocatorsListComplete(ctx, id, streamingpoliciesandstreaminglocators.DefaultStreamingLocatorsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StreamingPoliciesAndStreamingLocatorsClient.StreamingLocatorsListContentKeys`

```go
ctx := context.TODO()
id := streamingpoliciesandstreaminglocators.NewStreamingLocatorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingLocatorValue")

read, err := client.StreamingLocatorsListContentKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StreamingPoliciesAndStreamingLocatorsClient.StreamingLocatorsListPaths`

```go
ctx := context.TODO()
id := streamingpoliciesandstreaminglocators.NewStreamingLocatorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingLocatorValue")

read, err := client.StreamingLocatorsListPaths(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StreamingPoliciesAndStreamingLocatorsClient.StreamingPoliciesCreate`

```go
ctx := context.TODO()
id := streamingpoliciesandstreaminglocators.NewStreamingPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingPolicyValue")

payload := streamingpoliciesandstreaminglocators.StreamingPolicy{
	// ...
}


read, err := client.StreamingPoliciesCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StreamingPoliciesAndStreamingLocatorsClient.StreamingPoliciesDelete`

```go
ctx := context.TODO()
id := streamingpoliciesandstreaminglocators.NewStreamingPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingPolicyValue")

read, err := client.StreamingPoliciesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StreamingPoliciesAndStreamingLocatorsClient.StreamingPoliciesGet`

```go
ctx := context.TODO()
id := streamingpoliciesandstreaminglocators.NewStreamingPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingPolicyValue")

read, err := client.StreamingPoliciesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StreamingPoliciesAndStreamingLocatorsClient.StreamingPoliciesList`

```go
ctx := context.TODO()
id := streamingpoliciesandstreaminglocators.NewMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

// alternatively `client.StreamingPoliciesList(ctx, id, streamingpoliciesandstreaminglocators.DefaultStreamingPoliciesListOperationOptions())` can be used to do batched pagination
items, err := client.StreamingPoliciesListComplete(ctx, id, streamingpoliciesandstreaminglocators.DefaultStreamingPoliciesListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
