
## `github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/dpscertificate` Documentation

The `dpscertificate` SDK allows for interaction with Azure Resource Manager `deviceprovisioningservices` (API Version `2022-02-05`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/dpscertificate"
```


### Client Initialization

```go
client := dpscertificate.NewDpsCertificateClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DpsCertificateClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dpscertificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName", "certificateName")

payload := dpscertificate.CertificateResponse{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, dpscertificate.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DpsCertificateClient.Delete`

```go
ctx := context.TODO()
id := dpscertificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName", "certificateName")

read, err := client.Delete(ctx, id, dpscertificate.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DpsCertificateClient.GenerateVerificationCode`

```go
ctx := context.TODO()
id := dpscertificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName", "certificateName")

read, err := client.GenerateVerificationCode(ctx, id, dpscertificate.DefaultGenerateVerificationCodeOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DpsCertificateClient.Get`

```go
ctx := context.TODO()
id := dpscertificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName", "certificateName")

read, err := client.Get(ctx, id, dpscertificate.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DpsCertificateClient.List`

```go
ctx := context.TODO()
id := commonids.NewProvisioningServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DpsCertificateClient.VerifyCertificate`

```go
ctx := context.TODO()
id := dpscertificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "provisioningServiceName", "certificateName")

payload := dpscertificate.VerificationCodeRequest{
	// ...
}


read, err := client.VerifyCertificate(ctx, id, payload, dpscertificate.DefaultVerifyCertificateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
