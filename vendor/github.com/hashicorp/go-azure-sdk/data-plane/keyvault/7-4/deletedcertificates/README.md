
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/deletedcertificates` Documentation

The `deletedcertificates` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `7.4`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/deletedcertificates"
```


### Client Initialization

```go
client := deletedcertificates.NewDeletedCertificatesClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeletedCertificatesClient.GetDeletedCertificate`

```go
ctx := context.TODO()
id := deletedcertificates.NewDeletedcertificateID("deletedcertificateName")

read, err := client.GetDeletedCertificate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeletedCertificatesClient.GetDeletedCertificates`

```go
ctx := context.TODO()


// alternatively `client.GetDeletedCertificates(ctx, deletedcertificates.DefaultGetDeletedCertificatesOperationOptions())` can be used to do batched pagination
items, err := client.GetDeletedCertificatesComplete(ctx, deletedcertificates.DefaultGetDeletedCertificatesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeletedCertificatesClient.PurgeDeletedCertificate`

```go
ctx := context.TODO()
id := deletedcertificates.NewDeletedcertificateID("deletedcertificateName")

read, err := client.PurgeDeletedCertificate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeletedCertificatesClient.RecoverDeletedCertificate`

```go
ctx := context.TODO()
id := deletedcertificates.NewDeletedcertificateID("deletedcertificateName")

read, err := client.RecoverDeletedCertificate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
