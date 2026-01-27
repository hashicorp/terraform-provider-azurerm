
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/certificates` Documentation

The `certificates` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/certificates"
```


### Client Initialization

```go
client := certificates.NewCertificatesClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `CertificatesClient.BackupCertificate`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("certificateName")

read, err := client.BackupCertificate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.CreateCertificate`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("certificateName")

payload := certificates.CertificateCreateParameters{
	// ...
}


read, err := client.CreateCertificate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.DeleteCertificate`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("certificateName")

read, err := client.DeleteCertificate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.DeleteCertificateContacts`

```go
ctx := context.TODO()


read, err := client.DeleteCertificateContacts(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.DeleteCertificateIssuer`

```go
ctx := context.TODO()
id := certificates.NewIssuernameID("https://endpoint_url")

read, err := client.DeleteCertificateIssuer(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.DeleteCertificateOperation`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("certificateName")

read, err := client.DeleteCertificateOperation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.GetCertificate`

```go
ctx := context.TODO()
id := certificates.NewCertificateversionID("https://endpoint_url", "certificateName")

read, err := client.GetCertificate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.GetCertificateContacts`

```go
ctx := context.TODO()


read, err := client.GetCertificateContacts(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.GetCertificateIssuer`

```go
ctx := context.TODO()
id := certificates.NewIssuernameID("https://endpoint_url")

read, err := client.GetCertificateIssuer(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.GetCertificateIssuers`

```go
ctx := context.TODO()


// alternatively `client.GetCertificateIssuers(ctx, certificates.DefaultGetCertificateIssuersOperationOptions())` can be used to do batched pagination
items, err := client.GetCertificateIssuersComplete(ctx, certificates.DefaultGetCertificateIssuersOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CertificatesClient.GetCertificateOperation`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("certificateName")

read, err := client.GetCertificateOperation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.GetCertificatePolicy`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("certificateName")

read, err := client.GetCertificatePolicy(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.GetCertificateVersions`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("certificateName")

// alternatively `client.GetCertificateVersions(ctx, id, certificates.DefaultGetCertificateVersionsOperationOptions())` can be used to do batched pagination
items, err := client.GetCertificateVersionsComplete(ctx, id, certificates.DefaultGetCertificateVersionsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CertificatesClient.GetCertificates`

```go
ctx := context.TODO()


// alternatively `client.GetCertificates(ctx, certificates.DefaultGetCertificatesOperationOptions())` can be used to do batched pagination
items, err := client.GetCertificatesComplete(ctx, certificates.DefaultGetCertificatesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CertificatesClient.GetDeletedCertificate`

```go
ctx := context.TODO()
id := certificates.NewDeletedcertificateID("deletedcertificateName")

read, err := client.GetDeletedCertificate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.GetDeletedCertificates`

```go
ctx := context.TODO()


// alternatively `client.GetDeletedCertificates(ctx, certificates.DefaultGetDeletedCertificatesOperationOptions())` can be used to do batched pagination
items, err := client.GetDeletedCertificatesComplete(ctx, certificates.DefaultGetDeletedCertificatesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CertificatesClient.ImportCertificate`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("certificateName")

payload := certificates.CertificateImportParameters{
	// ...
}


read, err := client.ImportCertificate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.MergeCertificate`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("certificateName")

payload := certificates.CertificateMergeParameters{
	// ...
}


read, err := client.MergeCertificate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.PurgeDeletedCertificate`

```go
ctx := context.TODO()
id := certificates.NewDeletedcertificateID("deletedcertificateName")

read, err := client.PurgeDeletedCertificate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.RecoverDeletedCertificate`

```go
ctx := context.TODO()
id := certificates.NewDeletedcertificateID("deletedcertificateName")

read, err := client.RecoverDeletedCertificate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.RestoreCertificate`

```go
ctx := context.TODO()

payload := certificates.CertificateRestoreParameters{
	// ...
}


read, err := client.RestoreCertificate(ctx, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.SetCertificateContacts`

```go
ctx := context.TODO()

payload := certificates.Contacts{
	// ...
}


read, err := client.SetCertificateContacts(ctx, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.SetCertificateIssuer`

```go
ctx := context.TODO()
id := certificates.NewIssuernameID("https://endpoint_url")

payload := certificates.CertificateIssuerSetParameters{
	// ...
}


read, err := client.SetCertificateIssuer(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.UpdateCertificate`

```go
ctx := context.TODO()
id := certificates.NewCertificateversionID("https://endpoint_url", "certificateName")

payload := certificates.CertificateUpdateParameters{
	// ...
}


read, err := client.UpdateCertificate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.UpdateCertificateIssuer`

```go
ctx := context.TODO()
id := certificates.NewIssuernameID("https://endpoint_url")

payload := certificates.CertificateIssuerUpdateParameters{
	// ...
}


read, err := client.UpdateCertificateIssuer(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.UpdateCertificateOperation`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("certificateName")

payload := certificates.CertificateOperationUpdateParameter{
	// ...
}


read, err := client.UpdateCertificateOperation(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificatesClient.UpdateCertificatePolicy`

```go
ctx := context.TODO()
id := certificates.NewCertificateID("certificateName")

payload := certificates.CertificatePolicy{
	// ...
}


read, err := client.UpdateCertificatePolicy(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
