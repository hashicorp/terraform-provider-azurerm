
## `github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/certificateprofiles` Documentation

The `certificateprofiles` SDK allows for interaction with Azure Resource Manager `codesigning` (API Version `2025-10-13`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/certificateprofiles"
```


### Client Initialization

```go
client := certificateprofiles.NewCertificateProfilesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CertificateProfilesClient.Create`

```go
ctx := context.TODO()
id := certificateprofiles.NewCertificateProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "codeSigningAccountName", "certificateProfileName")

payload := certificateprofiles.CertificateProfile{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CertificateProfilesClient.Delete`

```go
ctx := context.TODO()
id := certificateprofiles.NewCertificateProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "codeSigningAccountName", "certificateProfileName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CertificateProfilesClient.Get`

```go
ctx := context.TODO()
id := certificateprofiles.NewCertificateProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "codeSigningAccountName", "certificateProfileName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificateProfilesClient.ListByCodeSigningAccount`

```go
ctx := context.TODO()
id := certificateprofiles.NewCodeSigningAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "codeSigningAccountName")

// alternatively `client.ListByCodeSigningAccount(ctx, id)` can be used to do batched pagination
items, err := client.ListByCodeSigningAccountComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CertificateProfilesClient.RevokeCertificate`

```go
ctx := context.TODO()
id := certificateprofiles.NewCertificateProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "codeSigningAccountName", "certificateProfileName")

payload := certificateprofiles.RevokeCertificate{
	// ...
}


read, err := client.RevokeCertificate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
