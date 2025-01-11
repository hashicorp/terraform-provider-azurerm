
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/certificate` Documentation

The `certificate` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/certificate"
```


### Client Initialization

```go
client := certificate.NewCertificateClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CertificateClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := certificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "certificateName")

payload := certificate.CertificateCreateOrUpdateParameters{
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


### Example Usage: `CertificateClient.Delete`

```go
ctx := context.TODO()
id := certificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "certificateName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificateClient.Get`

```go
ctx := context.TODO()
id := certificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "certificateName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificateClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := certificate.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

// alternatively `client.ListByAutomationAccount(ctx, id)` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id)
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
id := certificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "certificateName")

payload := certificate.CertificateUpdateParameters{
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
