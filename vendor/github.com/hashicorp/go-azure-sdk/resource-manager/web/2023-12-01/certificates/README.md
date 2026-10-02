
## `github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates` Documentation

The `certificates` SDK allows for interaction with Azure Resource Manager `web` (API Version `2023-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
```


### Client Initialization

```go
client := certificates.NewCertificatesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CertificatesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateName")

payload := certificates.Certificate{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CertificatesClient.Delete`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.Get`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id, certificates.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, certificates.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CertificatesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CertificatesClient.Update`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateName")

payload := certificates.CertificatePatchResource{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
