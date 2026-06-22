
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/hsmsecuritydomain` Documentation

The `hsmsecuritydomain` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `7.4`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/hsmsecuritydomain"
```


### Client Initialization

```go
client := hsmsecuritydomain.NewHSMSecurityDomainClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `HSMSecurityDomainClient.Download`

```go
ctx := context.TODO()

payload := hsmsecuritydomain.CertificateInfoObject{
	// ...
}


if err := client.DownloadThenPoll(ctx, payload); err != nil {
	// handle the error
}
```


### Example Usage: `HSMSecurityDomainClient.DownloadPending`

```go
ctx := context.TODO()


read, err := client.DownloadPending(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HSMSecurityDomainClient.TransferKey`

```go
ctx := context.TODO()


read, err := client.TransferKey(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HSMSecurityDomainClient.Upload`

```go
ctx := context.TODO()

payload := hsmsecuritydomain.SecurityDomainObject{
	// ...
}


if err := client.UploadThenPoll(ctx, payload); err != nil {
	// handle the error
}
```


### Example Usage: `HSMSecurityDomainClient.UploadPending`

```go
ctx := context.TODO()


read, err := client.UploadPending(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
