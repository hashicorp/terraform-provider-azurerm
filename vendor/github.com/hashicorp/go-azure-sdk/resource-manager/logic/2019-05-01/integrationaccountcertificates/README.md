
## `github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountcertificates` Documentation

The `integrationaccountcertificates` SDK allows for interaction with Azure Resource Manager `logic` (API Version `2019-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountcertificates"
```


### Client Initialization

```go
client := integrationaccountcertificates.NewIntegrationAccountCertificatesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IntegrationAccountCertificatesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := integrationaccountcertificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "certificateName")

payload := integrationaccountcertificates.IntegrationAccountCertificate{
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


### Example Usage: `IntegrationAccountCertificatesClient.Delete`

```go
ctx := context.TODO()
id := integrationaccountcertificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "certificateName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountCertificatesClient.Get`

```go
ctx := context.TODO()
id := integrationaccountcertificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "certificateName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountCertificatesClient.List`

```go
ctx := context.TODO()
id := integrationaccountcertificates.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName")

// alternatively `client.List(ctx, id, integrationaccountcertificates.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, integrationaccountcertificates.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
