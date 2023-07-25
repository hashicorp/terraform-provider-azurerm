
## `github.com/hashicorp/go-azure-sdk/resource-manager/videoanalyzer/2021-05-01-preview/videoanalyzers` Documentation

The `videoanalyzers` SDK allows for interaction with the Azure Resource Manager Service `videoanalyzer` (API Version `2021-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/videoanalyzer/2021-05-01-preview/videoanalyzers"
```


### Client Initialization

```go
client := videoanalyzers.NewVideoAnalyzersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VideoAnalyzersClient.LocationsCheckNameAvailability`

```go
ctx := context.TODO()
id := videoanalyzers.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

payload := videoanalyzers.CheckNameAvailabilityRequest{
	// ...
}


read, err := client.LocationsCheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzersClient.VideoAnalyzersCreateOrUpdate`

```go
ctx := context.TODO()
id := videoanalyzers.NewVideoAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "videoAnalyzerValue")

payload := videoanalyzers.VideoAnalyzer{
	// ...
}


read, err := client.VideoAnalyzersCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzersClient.VideoAnalyzersDelete`

```go
ctx := context.TODO()
id := videoanalyzers.NewVideoAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "videoAnalyzerValue")

read, err := client.VideoAnalyzersDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzersClient.VideoAnalyzersGet`

```go
ctx := context.TODO()
id := videoanalyzers.NewVideoAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "videoAnalyzerValue")

read, err := client.VideoAnalyzersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzersClient.VideoAnalyzersList`

```go
ctx := context.TODO()
id := videoanalyzers.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.VideoAnalyzersList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzersClient.VideoAnalyzersListBySubscription`

```go
ctx := context.TODO()
id := videoanalyzers.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.VideoAnalyzersListBySubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzersClient.VideoAnalyzersSyncStorageKeys`

```go
ctx := context.TODO()
id := videoanalyzers.NewVideoAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "videoAnalyzerValue")

payload := videoanalyzers.SyncStorageKeysInput{
	// ...
}


read, err := client.VideoAnalyzersSyncStorageKeys(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzersClient.VideoAnalyzersUpdate`

```go
ctx := context.TODO()
id := videoanalyzers.NewVideoAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "videoAnalyzerValue")

payload := videoanalyzers.VideoAnalyzerUpdate{
	// ...
}


read, err := client.VideoAnalyzersUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
