
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/certificates` Documentation

The `certificates` SDK allows for interaction with Azure Resource Manager `containerapps` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/certificates"
```


### Client Initialization

```go
client := certificates.NewCertificatesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CertificatesClient.ConnectedEnvironmentsCertificatesCreateOrUpdate`

```go
ctx := context.TODO()
id := certificates.NewConnectedEnvironmentCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectedEnvironmentName", "certificateName")

payload := certificates.Certificate{
	// ...
}


read, err := client.ConnectedEnvironmentsCertificatesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.ConnectedEnvironmentsCertificatesDelete`

```go
ctx := context.TODO()
id := certificates.NewConnectedEnvironmentCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectedEnvironmentName", "certificateName")

read, err := client.ConnectedEnvironmentsCertificatesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.ConnectedEnvironmentsCertificatesGet`

```go
ctx := context.TODO()
id := certificates.NewConnectedEnvironmentCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectedEnvironmentName", "certificateName")

read, err := client.ConnectedEnvironmentsCertificatesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.ConnectedEnvironmentsCertificatesList`

```go
ctx := context.TODO()
id := certificates.NewConnectedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectedEnvironmentName")

// alternatively `client.ConnectedEnvironmentsCertificatesList(ctx, id)` can be used to do batched pagination
items, err := client.ConnectedEnvironmentsCertificatesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CertificatesClient.ConnectedEnvironmentsCertificatesUpdate`

```go
ctx := context.TODO()
id := certificates.NewConnectedEnvironmentCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectedEnvironmentName", "certificateName")

payload := certificates.CertificatePatch{
	// ...
}


read, err := client.ConnectedEnvironmentsCertificatesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName", "certificateName")

payload := certificates.Certificate{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.Delete`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName", "certificateName")

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
id := certificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName", "certificateName")

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
id := certificates.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
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
id := certificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName", "certificateName")

payload := certificates.CertificatePatch{
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
