
## `github.com/hashicorp/go-azure-sdk/resource-manager/batch/2023-05-01/certificate` Documentation

The `certificate` SDK allows for interaction with the Azure Resource Manager Service `batch` (API Version `2023-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/batch/2023-05-01/certificate"
```


### Client Initialization

```go
client := certificate.NewCertificateClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CertificateClient.CancelDeletion`

```go
ctx := context.TODO()
id := certificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountValue", "certificateValue")

read, err := client.CancelDeletion(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificateClient.Create`

```go
ctx := context.TODO()
id := certificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountValue", "certificateValue")

payload := certificate.CertificateCreateOrUpdateParameters{
	// ...
}


read, err := client.Create(ctx, id, payload, certificate.DefaultCreateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificateClient.Delete`

```go
ctx := context.TODO()
id := certificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountValue", "certificateValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CertificateClient.Get`

```go
ctx := context.TODO()
id := certificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountValue", "certificateValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificateClient.ListByBatchAccount`

```go
ctx := context.TODO()
id := certificate.NewBatchAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountValue")

// alternatively `client.ListByBatchAccount(ctx, id, certificate.DefaultListByBatchAccountOperationOptions())` can be used to do batched pagination
items, err := client.ListByBatchAccountComplete(ctx, id, certificate.DefaultListByBatchAccountOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CertificateClient.Update`

```go
ctx := context.TODO()
id := certificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountValue", "certificateValue")

payload := certificate.CertificateCreateOrUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload, certificate.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
