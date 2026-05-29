
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/gallerysharingupdate` Documentation

The `gallerysharingupdate` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2022-03-03`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/gallerysharingupdate"
```


### Client Initialization

```go
client := gallerysharingupdate.NewGallerySharingUpdateClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GallerySharingUpdateClient.GallerySharingProfileUpdate`

```go
ctx := context.TODO()
id := commonids.NewSharedImageGalleryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryName")

payload := gallerysharingupdate.SharingUpdate{
	// ...
}


if err := client.GallerySharingProfileUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
