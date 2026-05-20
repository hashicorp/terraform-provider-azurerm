
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachineimages` Documentation

The `virtualmachineimages` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachineimages"
```


### Client Initialization

```go
client := virtualmachineimages.NewVirtualMachineImagesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualMachineImagesClient.EdgeZoneGet`

```go
ctx := context.TODO()
id := virtualmachineimages.NewOfferSkuVersionID("12345678-1234-9876-4563-123456789012", "locationName", "edgeZoneName", "publisherName", "offerName", "skuName", "versionName")

read, err := client.EdgeZoneGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineImagesClient.EdgeZoneList`

```go
ctx := context.TODO()
id := virtualmachineimages.NewOfferSkuID("12345678-1234-9876-4563-123456789012", "locationName", "edgeZoneName", "publisherName", "offerName", "skuName")

read, err := client.EdgeZoneList(ctx, id, virtualmachineimages.DefaultEdgeZoneListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineImagesClient.EdgeZoneListOffers`

```go
ctx := context.TODO()
id := virtualmachineimages.NewEdgeZonePublisherID("12345678-1234-9876-4563-123456789012", "locationName", "edgeZoneName", "publisherName")

read, err := client.EdgeZoneListOffers(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineImagesClient.EdgeZoneListPublishers`

```go
ctx := context.TODO()
id := virtualmachineimages.NewEdgeZoneID("12345678-1234-9876-4563-123456789012", "locationName", "edgeZoneName")

read, err := client.EdgeZoneListPublishers(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineImagesClient.EdgeZoneListSkus`

```go
ctx := context.TODO()
id := virtualmachineimages.NewVMImageOfferID("12345678-1234-9876-4563-123456789012", "locationName", "edgeZoneName", "publisherName", "offerName")

read, err := client.EdgeZoneListSkus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineImagesClient.Get`

```go
ctx := context.TODO()
id := virtualmachineimages.NewSkuVersionID("12345678-1234-9876-4563-123456789012", "locationName", "publisherName", "offerName", "skuName", "versionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineImagesClient.List`

```go
ctx := context.TODO()
id := virtualmachineimages.NewSkuID("12345678-1234-9876-4563-123456789012", "locationName", "publisherName", "offerName", "skuName")

read, err := client.List(ctx, id, virtualmachineimages.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineImagesClient.ListByEdgeZone`

```go
ctx := context.TODO()
id := virtualmachineimages.NewEdgeZoneID("12345678-1234-9876-4563-123456789012", "locationName", "edgeZoneName")

// alternatively `client.ListByEdgeZone(ctx, id)` can be used to do batched pagination
items, err := client.ListByEdgeZoneComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualMachineImagesClient.ListOffers`

```go
ctx := context.TODO()
id := virtualmachineimages.NewPublisherID("12345678-1234-9876-4563-123456789012", "locationName", "publisherName")

read, err := client.ListOffers(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineImagesClient.ListPublishers`

```go
ctx := context.TODO()
id := virtualmachineimages.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

read, err := client.ListPublishers(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualMachineImagesClient.ListSkus`

```go
ctx := context.TODO()
id := virtualmachineimages.NewOfferID("12345678-1234-9876-4563-123456789012", "locationName", "publisherName", "offerName")

read, err := client.ListSkus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
