
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectlocalrulestackresources` Documentation

The `certificateobjectlocalrulestackresources` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2025-10-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectlocalrulestackresources"
```


### Client Initialization

```go
client := certificateobjectlocalrulestackresources.NewCertificateObjectLocalRulestackResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CertificateObjectLocalRulestackResourcesClient.CertificateObjectLocalRulestackCreateOrUpdate`

```go
ctx := context.TODO()
id := certificateobjectlocalrulestackresources.NewLocalRulestackCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "certificateName")

payload := certificateobjectlocalrulestackresources.CertificateObjectLocalRulestackResource{
	// ...
}


if err := client.CertificateObjectLocalRulestackCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CertificateObjectLocalRulestackResourcesClient.CertificateObjectLocalRulestackDelete`

```go
ctx := context.TODO()
id := certificateobjectlocalrulestackresources.NewLocalRulestackCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "certificateName")

if err := client.CertificateObjectLocalRulestackDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CertificateObjectLocalRulestackResourcesClient.CertificateObjectLocalRulestackGet`

```go
ctx := context.TODO()
id := certificateobjectlocalrulestackresources.NewLocalRulestackCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName", "certificateName")

read, err := client.CertificateObjectLocalRulestackGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificateObjectLocalRulestackResourcesClient.CertificateObjectLocalRulestackListByLocalRulestacks`

```go
ctx := context.TODO()
id := certificateobjectlocalrulestackresources.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.CertificateObjectLocalRulestackListByLocalRulestacks(ctx, id)` can be used to do batched pagination
items, err := client.CertificateObjectLocalRulestackListByLocalRulestacksComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
