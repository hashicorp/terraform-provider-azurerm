
## `github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/assetsandassetfilters` Documentation

The `assetsandassetfilters` SDK allows for interaction with the Azure Resource Manager Service `media` (API Version `2021-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/assetsandassetfilters"
```


### Client Initialization

```go
client := assetsandassetfilters.NewAssetsAndAssetFiltersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AssetsAndAssetFiltersClient.AssetFiltersCreateOrUpdate`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewAssetFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue", "assetFilterValue")

payload := assetsandassetfilters.AssetFilter{
	// ...
}


read, err := client.AssetFiltersCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssetsAndAssetFiltersClient.AssetFiltersDelete`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewAssetFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue", "assetFilterValue")

read, err := client.AssetFiltersDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssetsAndAssetFiltersClient.AssetFiltersGet`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewAssetFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue", "assetFilterValue")

read, err := client.AssetFiltersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssetsAndAssetFiltersClient.AssetFiltersList`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewAssetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue")

// alternatively `client.AssetFiltersList(ctx, id)` can be used to do batched pagination
items, err := client.AssetFiltersListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AssetsAndAssetFiltersClient.AssetFiltersUpdate`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewAssetFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue", "assetFilterValue")

payload := assetsandassetfilters.AssetFilter{
	// ...
}


read, err := client.AssetFiltersUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssetsAndAssetFiltersClient.AssetsCreateOrUpdate`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewAssetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue")

payload := assetsandassetfilters.Asset{
	// ...
}


read, err := client.AssetsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssetsAndAssetFiltersClient.AssetsDelete`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewAssetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue")

read, err := client.AssetsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssetsAndAssetFiltersClient.AssetsGet`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewAssetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue")

read, err := client.AssetsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssetsAndAssetFiltersClient.AssetsGetEncryptionKey`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewAssetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue")

read, err := client.AssetsGetEncryptionKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssetsAndAssetFiltersClient.AssetsList`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

// alternatively `client.AssetsList(ctx, id, assetsandassetfilters.DefaultAssetsListOperationOptions())` can be used to do batched pagination
items, err := client.AssetsListComplete(ctx, id, assetsandassetfilters.DefaultAssetsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AssetsAndAssetFiltersClient.AssetsListContainerSas`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewAssetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue")

payload := assetsandassetfilters.ListContainerSasInput{
	// ...
}


read, err := client.AssetsListContainerSas(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssetsAndAssetFiltersClient.AssetsListStreamingLocators`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewAssetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue")

read, err := client.AssetsListStreamingLocators(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssetsAndAssetFiltersClient.AssetsUpdate`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewAssetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue")

payload := assetsandassetfilters.Asset{
	// ...
}


read, err := client.AssetsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssetsAndAssetFiltersClient.TracksCreateOrUpdate`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewTrackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue", "trackValue")

payload := assetsandassetfilters.AssetTrack{
	// ...
}


if err := client.TracksCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AssetsAndAssetFiltersClient.TracksDelete`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewTrackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue", "trackValue")

if err := client.TracksDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AssetsAndAssetFiltersClient.TracksGet`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewTrackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue", "trackValue")

read, err := client.TracksGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssetsAndAssetFiltersClient.TracksList`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewAssetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue")

read, err := client.TracksList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssetsAndAssetFiltersClient.TracksUpdate`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewTrackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue", "trackValue")

payload := assetsandassetfilters.AssetTrack{
	// ...
}


if err := client.TracksUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AssetsAndAssetFiltersClient.TracksUpdateTrackData`

```go
ctx := context.TODO()
id := assetsandassetfilters.NewTrackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "assetValue", "trackValue")

if err := client.TracksUpdateTrackDataThenPoll(ctx, id); err != nil {
	// handle the error
}
```
