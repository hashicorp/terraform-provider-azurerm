
## `github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/certificates` Documentation

The `certificates` SDK allows for interaction with Azure Resource Manager `batch` (API Version `2024-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/certificates"
```


### Client Initialization

```go
client := certificates.NewCertificatesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CertificatesClient.CertificateCancelDeletion`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "certificateName")

read, err := client.CertificateCancelDeletion(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.CertificateCreate`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "certificateName")

payload := certificates.CertificateCreateOrUpdateParameters{
	// ...
}


read, err := client.CertificateCreate(ctx, id, payload, certificates.DefaultCertificateCreateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.CertificateDelete`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "certificateName")

if err := client.CertificateDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CertificatesClient.CertificateGet`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "certificateName")

read, err := client.CertificateGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.CertificateListByBatchAccount`

```go
ctx := context.TODO()
id := certificates.NewBatchAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName")

// alternatively `client.CertificateListByBatchAccount(ctx, id, certificates.DefaultCertificateListByBatchAccountOperationOptions())` can be used to do batched pagination
items, err := client.CertificateListByBatchAccountComplete(ctx, id, certificates.DefaultCertificateListByBatchAccountOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CertificatesClient.CertificateUpdate`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "certificateName")

payload := certificates.CertificateCreateOrUpdateParameters{
	// ...
}


read, err := client.CertificateUpdate(ctx, id, payload, certificates.DefaultCertificateUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
