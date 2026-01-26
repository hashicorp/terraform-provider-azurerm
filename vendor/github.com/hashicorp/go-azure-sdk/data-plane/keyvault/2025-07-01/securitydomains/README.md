
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/securitydomains` Documentation

The `securitydomains` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/securitydomains"
```


### Client Initialization

```go
client := securitydomains.NewSecuritydomainsClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `SecuritydomainsClient.HSMSecurityDomainDownload`

```go
ctx := context.TODO()

payload := securitydomains.CertificateInfoObject{
	// ...
}


if err := client.HSMSecurityDomainDownloadThenPoll(ctx, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SecuritydomainsClient.HSMSecurityDomainDownloadPending`

```go
ctx := context.TODO()


read, err := client.HSMSecurityDomainDownloadPending(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecuritydomainsClient.HSMSecurityDomainTransferKey`

```go
ctx := context.TODO()


read, err := client.HSMSecurityDomainTransferKey(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecuritydomainsClient.HSMSecurityDomainUpload`

```go
ctx := context.TODO()

payload := securitydomains.SecurityDomainObject{
	// ...
}


if err := client.HSMSecurityDomainUploadThenPoll(ctx, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SecuritydomainsClient.HSMSecurityDomainUploadPending`

```go
ctx := context.TODO()


read, err := client.HSMSecurityDomainUploadPending(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
