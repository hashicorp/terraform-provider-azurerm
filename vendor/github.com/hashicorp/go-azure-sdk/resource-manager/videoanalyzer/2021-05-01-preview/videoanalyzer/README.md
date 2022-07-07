
## `github.com/hashicorp/go-azure-sdk/resource-manager/videoanalyzer/2021-05-01-preview/videoanalyzer` Documentation

The `videoanalyzer` SDK allows for interaction with the Azure Resource Manager Service `videoanalyzer` (API Version `2021-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/videoanalyzer/2021-05-01-preview/videoanalyzer"
```


### Client Initialization

```go
client := videoanalyzer.NewVideoAnalyzerClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VideoAnalyzerClient.AccessPoliciesCreateOrUpdate`

```go
ctx := context.TODO()
id := videoanalyzer.NewAccessPoliciesID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "accessPolicyValue")

payload := videoanalyzer.AccessPolicyEntity{
	// ...
}


read, err := client.AccessPoliciesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.AccessPoliciesDelete`

```go
ctx := context.TODO()
id := videoanalyzer.NewAccessPoliciesID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "accessPolicyValue")

read, err := client.AccessPoliciesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.AccessPoliciesGet`

```go
ctx := context.TODO()
id := videoanalyzer.NewAccessPoliciesID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "accessPolicyValue")

read, err := client.AccessPoliciesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.AccessPoliciesList`

```go
ctx := context.TODO()
id := videoanalyzer.NewVideoAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

// alternatively `client.AccessPoliciesList(ctx, id, videoanalyzer.DefaultAccessPoliciesListOperationOptions())` can be used to do batched pagination
items, err := client.AccessPoliciesListComplete(ctx, id, videoanalyzer.DefaultAccessPoliciesListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VideoAnalyzerClient.AccessPoliciesUpdate`

```go
ctx := context.TODO()
id := videoanalyzer.NewAccessPoliciesID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "accessPolicyValue")

payload := videoanalyzer.AccessPolicyEntity{
	// ...
}


read, err := client.AccessPoliciesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.EdgeModulesCreateOrUpdate`

```go
ctx := context.TODO()
id := videoanalyzer.NewEdgeModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "edgeModuleValue")

payload := videoanalyzer.EdgeModuleEntity{
	// ...
}


read, err := client.EdgeModulesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.EdgeModulesDelete`

```go
ctx := context.TODO()
id := videoanalyzer.NewEdgeModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "edgeModuleValue")

read, err := client.EdgeModulesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.EdgeModulesGet`

```go
ctx := context.TODO()
id := videoanalyzer.NewEdgeModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "edgeModuleValue")

read, err := client.EdgeModulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.EdgeModulesList`

```go
ctx := context.TODO()
id := videoanalyzer.NewVideoAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

// alternatively `client.EdgeModulesList(ctx, id, videoanalyzer.DefaultEdgeModulesListOperationOptions())` can be used to do batched pagination
items, err := client.EdgeModulesListComplete(ctx, id, videoanalyzer.DefaultEdgeModulesListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VideoAnalyzerClient.EdgeModulesListProvisioningToken`

```go
ctx := context.TODO()
id := videoanalyzer.NewEdgeModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "edgeModuleValue")

payload := videoanalyzer.ListProvisioningTokenInput{
	// ...
}


read, err := client.EdgeModulesListProvisioningToken(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.LocationsCheckNameAvailability`

```go
ctx := context.TODO()
id := videoanalyzer.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

payload := videoanalyzer.CheckNameAvailabilityRequest{
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


### Example Usage: `VideoAnalyzerClient.VideoAnalyzersCreateOrUpdate`

```go
ctx := context.TODO()
id := videoanalyzer.NewVideoAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := videoanalyzer.VideoAnalyzer{
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


### Example Usage: `VideoAnalyzerClient.VideoAnalyzersDelete`

```go
ctx := context.TODO()
id := videoanalyzer.NewVideoAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.VideoAnalyzersDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.VideoAnalyzersGet`

```go
ctx := context.TODO()
id := videoanalyzer.NewVideoAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.VideoAnalyzersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.VideoAnalyzersList`

```go
ctx := context.TODO()
id := videoanalyzer.NewResourceGroupID()

read, err := client.VideoAnalyzersList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.VideoAnalyzersListBySubscription`

```go
ctx := context.TODO()
id := videoanalyzer.NewSubscriptionID()

read, err := client.VideoAnalyzersListBySubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.VideoAnalyzersSyncStorageKeys`

```go
ctx := context.TODO()
id := videoanalyzer.NewVideoAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := videoanalyzer.SyncStorageKeysInput{
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


### Example Usage: `VideoAnalyzerClient.VideoAnalyzersUpdate`

```go
ctx := context.TODO()
id := videoanalyzer.NewVideoAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := videoanalyzer.VideoAnalyzerUpdate{
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


### Example Usage: `VideoAnalyzerClient.VideosCreateOrUpdate`

```go
ctx := context.TODO()
id := videoanalyzer.NewVideoID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "videoValue")

payload := videoanalyzer.VideoEntity{
	// ...
}


read, err := client.VideosCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.VideosDelete`

```go
ctx := context.TODO()
id := videoanalyzer.NewVideoID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "videoValue")

read, err := client.VideosDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.VideosGet`

```go
ctx := context.TODO()
id := videoanalyzer.NewVideoID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "videoValue")

read, err := client.VideosGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.VideosList`

```go
ctx := context.TODO()
id := videoanalyzer.NewVideoAnalyzerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

// alternatively `client.VideosList(ctx, id, videoanalyzer.DefaultVideosListOperationOptions())` can be used to do batched pagination
items, err := client.VideosListComplete(ctx, id, videoanalyzer.DefaultVideosListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VideoAnalyzerClient.VideosListStreamingToken`

```go
ctx := context.TODO()
id := videoanalyzer.NewVideoID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "videoValue")

read, err := client.VideosListStreamingToken(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VideoAnalyzerClient.VideosUpdate`

```go
ctx := context.TODO()
id := videoanalyzer.NewVideoID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "videoValue")

payload := videoanalyzer.VideoEntity{
	// ...
}


read, err := client.VideosUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
